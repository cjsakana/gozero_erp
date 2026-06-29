package model

import (
	"context"
	"database/sql"
	"erp/app/inventory/rpc/internal/types"
	"errors"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ InventoryModel = (*customInventoryModel)(nil)

type (
	// InventoryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customInventoryModel.
	InventoryModel interface {
		inventoryModel
		UpdateTransactCtx(ctx context.Context, newData *Inventory, adjustType int64, inventoryTransaction *InventoryTransaction) error
		XUpdate(ctx context.Context, newData *Inventory) error
		Search(ctx context.Context, data *types.SearchInventoryParams) ([]*Inventory, int64, error)
		LowStockAlert(ctx context.Context) ([]*Inventory, error)
		InsertTransactCtx(ctx context.Context, inventory *Inventory, inventoryTransaction *InventoryTransaction) error
		GetUsedCapacity(ctx context.Context, warehouseId int64) (float64, error)
	}

	customInventoryModel struct {
		*defaultInventoryModel
		conn sqlx.SqlConn
	}
)

// NewInventoryModel returns a model for the database table.
func NewInventoryModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) InventoryModel {
	return &customInventoryModel{
		defaultInventoryModel: newInventoryModel(conn, c, opts...),
		conn:                  conn,
	}
}

const (
	AdjustTypeIncrease = 1
	AdjustTypeDecrease = 2
	AdjustTypeSet      = 3
)

func (m *customInventoryModel) UpdateTransactCtx(ctx context.Context, newData *Inventory, adjustType int64, inventoryTransaction *InventoryTransaction) error {
	var setClauses []string
	var args []interface{}

	switch adjustType {
	case AdjustTypeIncrease:
		setClauses = append(setClauses, "current_stock = current_stock + ?")
		args = append(args, newData.CurrentStock)
	case AdjustTypeDecrease:
		setClauses = append(setClauses, "current_stock = current_stock - ?")
		args = append(args, newData.CurrentStock)
	case AdjustTypeSet:
		setClauses = append(setClauses, "current_stock = ?")
		args = append(args, newData.CurrentStock)
	}

	if newData.SafetyStock != 0 {
		setClauses = append(setClauses, "safety_stock = ?")
		args = append(args, newData.SafetyStock)
	}
	if newData.LockedStock != 0 {
		setClauses = append(setClauses, "locked_stock = ?")
		args = append(args, newData.LockedStock)
	}

	err := m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		query := fmt.Sprintf("update %s set %s where `inventory_id` = ?", m.table, strings.Join(setClauses, ", "))
		args = append(args, newData.InventoryId)
		_, err := session.ExecCtx(ctx, query, args...)
		if err != nil {
			return err
		}

		query2 := fmt.Sprintf("insert into inventory_transaction (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", inventoryTransactionRowsExpectAutoSet)
		_, err = session.ExecCtx(ctx, query2, inventoryTransaction.Id, inventoryTransaction.ProductId, inventoryTransaction.WarehouseId,
			inventoryTransaction.BatchId, inventoryTransaction.TransactionType, inventoryTransaction.Quantity,
			inventoryTransaction.ReferenceType, inventoryTransaction.ReferenceId, inventoryTransaction.OperatorId)
		if err != nil {
			return err
		}
		return nil
	})

	erpInventoryInventoryInventoryIdKey := fmt.Sprintf("%s%v", cacheErpInventoryInventoryInventoryIdPrefix, newData.InventoryId)
	m.CachedConn.DelCache(erpInventoryInventoryInventoryIdKey)
	return err
}

func (m *customInventoryModel) XUpdate(ctx context.Context, newData *Inventory) error {
	var setClauses []string
	var args []interface{}

	if newData.SafetyStock != 0 {
		setClauses = append(setClauses, "safety_stock = ?")
		args = append(args, newData.SafetyStock)
	}
	if newData.LockedStock != 0 {
		setClauses = append(setClauses, "locked_stock = ?")
		args = append(args, newData.LockedStock)
	}

	// 如果没有要更新的字段，直接返回
	if len(setClauses) == 0 {
		return nil
	}

	erpInventoryInventoryInventoryIdKey := fmt.Sprintf("%s%v", cacheErpInventoryInventoryInventoryIdPrefix, newData.InventoryId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `inventory_id` = ?", m.table, strings.Join(setClauses, ", "))
		args = append(args, newData.InventoryId)
		return conn.ExecCtx(ctx, query, args...)
	}, erpInventoryInventoryInventoryIdKey)
	return err
}

func (m *customInventoryModel) Search(ctx context.Context, data *types.SearchInventoryParams) ([]*Inventory, int64, error) {
	var inventories []*Inventory

	conditions := []string{}
	args := []any{}

	if data.ProductId != 0 {
		conditions = append(conditions, "product_id = ?")
		args = append(args, data.ProductId)
	}
	if data.WarehouseId != 0 {
		conditions = append(conditions, "warehouse_id = ?")
		args = append(args, data.WarehouseId)
	}
	if data.CurrentStock != 0 {
		conditions = append(conditions, "current_stock = ?")
		args = append(args, data.CurrentStock)
	}
	if data.SafetyStock != 0 {
		conditions = append(conditions, "safety_stock = ?")
		args = append(args, data.SafetyStock)
	}
	if data.LockedStock != 0 {
		conditions = append(conditions, "locked_stock = ?")
		args = append(args, data.LockedStock)
	}

	// 构建完整 SQL
	sql := fmt.Sprintf("select %s from %s", inventoryRows, m.table)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM  %s", m.table)
	if len(conditions) > 0 {
		where := " where " + strings.Join(conditions, " AND ")
		sql += where
		countQuery += where
	}

	sql += " ORDER BY product_id ASC "

	// 查询总数
	var total int64
	err := m.QueryRowNoCacheCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 添加分页
	if data.Limit != -1 { // 约定 -1 表示查询全部
		sql += fmt.Sprintf(" LIMIT %d OFFSET %d", data.Limit, (data.Page-1)*data.Limit)
	}
	err = m.QueryRowsNoCacheCtx(ctx, &inventories, sql, args...)

	switch {
	case err == nil:
		return inventories, total, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, 0, ErrNotFound
	default:
		return nil, 0, err
	}
}

func (m *customInventoryModel) LowStockAlert(ctx context.Context) ([]*Inventory, error) {
	var inventories []*Inventory
	query := fmt.Sprintf("SELECT * FROM %s WHERE current_stock <= safety_stock AND safety_stock > 0", m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &inventories, query)

	switch {
	case err == nil:
		return inventories, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customInventoryModel) InsertTransactCtx(ctx context.Context, inventory *Inventory, inventoryTransaction *InventoryTransaction) error {
	err := m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, inventoryRowsExpectAutoSet)

		_, err := session.ExecCtx(ctx, query, inventory.InventoryId, inventory.ProductId, inventory.WarehouseId, inventory.CurrentStock, inventory.SafetyStock, inventory.LockedStock)
		if err != nil {
			return err
		}

		query2 := fmt.Sprintf("insert into inventory_transaction (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", inventoryTransactionRowsExpectAutoSet)
		_, err = session.ExecCtx(ctx, query2, inventoryTransaction.Id, inventoryTransaction.ProductId, inventoryTransaction.WarehouseId,
			inventoryTransaction.BatchId, inventoryTransaction.TransactionType, inventoryTransaction.Quantity,
			inventoryTransaction.ReferenceType, inventoryTransaction.ReferenceId, inventoryTransaction.OperatorId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// GetUsedCapacity 获取仓库已使用容量
func (m *customInventoryModel) GetUsedCapacity(ctx context.Context, warehouseId int64) (float64, error) {
	var usedCapacity float64
	query := fmt.Sprintf("SELECT COALESCE(SUM(current_stock), 0) FROM %s WHERE warehouse_id = ?", m.table)
	err := m.QueryRowNoCacheCtx(ctx, &usedCapacity, query, warehouseId)
	if err != nil {
		return 0, err
	}
	return usedCapacity, nil
}

package model

import (
	"context"
	"erp/app/inventory/rpc/internal/types"
	"erp/common/xtime"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
)

var _ InventoryTransactionModel = (*customInventoryTransactionModel)(nil)

type (
	// InventoryTransactionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customInventoryTransactionModel.
	InventoryTransactionModel interface {
		inventoryTransactionModel
		Search(ctx context.Context, data *types.SearchInventoryTransactionParams) ([]*InventoryTransaction, int64, error)
	}

	customInventoryTransactionModel struct {
		*defaultInventoryTransactionModel
	}
)

// NewInventoryTransactionModel returns a model for the database table.
func NewInventoryTransactionModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) InventoryTransactionModel {
	return &customInventoryTransactionModel{
		defaultInventoryTransactionModel: newInventoryTransactionModel(conn, c, opts...),
	}
}

//func (m *customInventoryTransactionModel) XUpdate(ctx context.Context, newData *InventoryTransaction) error {
//	var setClauses []string
//	var args []interface{}
//
//	if newData.CurrentStock != 0 {
//		setClauses = append(setClauses, "current_stock = ?")
//		args = append(args, newData.CurrentStock)
//	}
//	if newData.SafetyStock != 0 {
//		setClauses = append(setClauses, "safety_stock = ?")
//		args = append(args, newData.SafetyStock)
//	}
//	if newData.LockedStock != 0 {
//		setClauses = append(setClauses, "locked_stock = ?")
//		args = append(args, newData.LockedStock)
//	}
//
//	erpInventoryInventoryTransactionIdKey := fmt.Sprintf("%s%v", cacheErpInventoryInventoryTransactionIdPrefix, newData.Id)
//	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
//		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, strings.Join(setClauses, ", "))
//		args = append(args, newData.Id)
//		return conn.ExecCtx(ctx, query, args...)
//	}, erpInventoryInventoryTransactionIdKey)
//	return err
//}

func (m *customInventoryTransactionModel) Search(ctx context.Context, data *types.SearchInventoryTransactionParams) ([]*InventoryTransaction, int64, error) {
	var records []*InventoryTransaction

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
	if data.BatchId != 0 {
		conditions = append(conditions, "batch_id = ?")
		args = append(args, data.BatchId)
	}
	if data.TransactionType != 0 {
		conditions = append(conditions, "transaction_type = ?")
		args = append(args, data.TransactionType)
	}
	if data.ReferenceType != 0 {
		conditions = append(conditions, "reference_type = ?")
		args = append(args, data.ReferenceType)
	}
	if !xtime.IsZeroTime(data.StartTime) {
		conditions = append(conditions, "created_at >= ?")
		args = append(args, data.StartTime)
	}
	if !xtime.IsZeroTime(data.EndTime) {
		conditions = append(conditions, "created_at <= ?")
		args = append(args, data.EndTime)
	}

	// 构建完整 SQL
	sql := fmt.Sprintf("select %s from %s", inventoryTransactionRows, m.table)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM  %s", m.table)
	if len(conditions) > 0 {
		where := " where " + strings.Join(conditions, " AND ")
		sql += where
		countQuery += where
	}

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
	err = m.QueryRowsNoCacheCtx(ctx, &records, sql, args...)

	switch {
	case err == nil:
		return records, total, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, 0, ErrNotFound
	default:
		return nil, 0, err
	}
}

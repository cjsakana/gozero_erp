package model

import (
	"context"
	"database/sql"
	"erp/app/inventory/rpc/internal/types"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
)

var _ WarehouseModel = (*customWarehouseModel)(nil)

type (
	// WarehouseModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWarehouseModel.
	WarehouseModel interface {
		warehouseModel
		FindOneByNo(ctx context.Context, no string) (*Warehouse, error)
		XUpdate(ctx context.Context, newData *Warehouse) error
		Search(ctx context.Context, data *types.SearchWarehouseParams) ([]*Warehouse, int64, error)
	}

	customWarehouseModel struct {
		*defaultWarehouseModel
	}
)

// NewWarehouseModel returns a model for the database table.
func NewWarehouseModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) WarehouseModel {
	return &customWarehouseModel{
		defaultWarehouseModel: newWarehouseModel(conn, c, opts...),
	}
}
func (m *customWarehouseModel) XUpdate(ctx context.Context, newData *Warehouse) error {
	var setClauses []string
	var args []interface{}

	if newData.Name != "" {
		setClauses = append(setClauses, "name = ?")
		args = append(args, newData.Name)
	}
	if newData.Location.String != "" {
		setClauses = append(setClauses, "location = ?")
		args = append(args, newData.Location.String)
	}
	if newData.ManagerId.Valid {
		setClauses = append(setClauses, "manager_id = ?")
		args = append(args, newData.ManagerId.Int64)
	}
	if newData.Capacity.Float64 != 0 {
		setClauses = append(setClauses, "capacity = ?")
		args = append(args, newData.Capacity.Float64)
	}
	if newData.IsActive != 0 {
		setClauses = append(setClauses, "is_active = ?")
		args = append(args, newData.IsActive)
	}

	erpInventoryWarehouseIdKey := fmt.Sprintf("%s%v", cacheErpInventoryWarehouseIdPrefix, newData.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, strings.Join(setClauses, ", "))
		args = append(args, newData.Id)
		return conn.ExecCtx(ctx, query, args...)
	}, erpInventoryWarehouseIdKey)
	return err
}

func (m *customWarehouseModel) FindOneByNo(ctx context.Context, no string) (*Warehouse, error) {
	var resp Warehouse
	query := fmt.Sprintf("select %s from %s where `no` = ? limit 1", warehouseRows, m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, no)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customWarehouseModel) Search(ctx context.Context, data *types.SearchWarehouseParams) ([]*Warehouse, int64, error) {
	var records []*Warehouse

	conditions := []string{}
	args := []any{}

	// 处理 OR 条件（name 或 location）
	if data.Name != "" || data.Location != "" {
		orConditions := []string{}

		if data.Name != "" {
			orConditions = append(orConditions, "name like ?")
			args = append(args, "%"+data.Name+"%")
		}
		if data.Location != "" {
			orConditions = append(orConditions, "location like ?")
			args = append(args, "%"+data.Location+"%")
		}

		if len(orConditions) > 0 {
			orClause := "(" + strings.Join(orConditions, " OR ") + ")"
			conditions = append(conditions, orClause)
		}
	}
	if data.IsActive != 0 {
		conditions = append(conditions, "is_active = ?")
		args = append(args, data.IsActive)
	}

	// 构建完整 SQL
	sql := fmt.Sprintf("select %s from %s", warehouseRows, m.table)
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

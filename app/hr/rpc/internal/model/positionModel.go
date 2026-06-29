package model

import (
	"context"
	"database/sql"
	"erp/app/hr/rpc/internal/types"
	"errors"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PositionModel = (*customPositionModel)(nil)

type (
	// PositionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPositionModel.
	PositionModel interface {
		positionModel
		XUpdate(ctx context.Context, newData *Position) error
		Search(ctx context.Context, data *types.SearchPositionParams) ([]*Position, int64, error)
	}

	customPositionModel struct {
		*defaultPositionModel
		conn sqlx.SqlConn
	}
)

// NewPositionModel returns a model for the database table.
func NewPositionModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PositionModel {
	return &customPositionModel{
		defaultPositionModel: newPositionModel(conn, c, opts...),
		conn:                 conn,
	}
}

func (m *customPositionModel) XUpdate(ctx context.Context, newData *Position) error {
	var setClauses []string
	var args []interface{}

	if newData.Name != "" {
		setClauses = append(setClauses, "name = ?")
		args = append(args, newData.Name)
	}
	if newData.Description.String != "" {
		setClauses = append(setClauses, "description = ?")
		args = append(args, newData.Description.String)
	}

	erpHrPositionIdKey := fmt.Sprintf("%s%v", cacheErpHrPositionIdPrefix, newData.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, strings.Join(setClauses, ", "))
		args = append(args, newData.Id)
		return conn.ExecCtx(ctx, query, args...)
	}, erpHrPositionIdKey)
	return err
}

func (m *customPositionModel) Search(ctx context.Context, data *types.SearchPositionParams) ([]*Position, int64, error) {
	var position []*Position

	conditions := []string{}
	args := []any{}

	if data.Name != "" {
		conditions = append(conditions, "name = ?")
		args = append(args, data.Name)
	}

	if data.Description != "" {
		conditions = append(conditions, "description LIKE ?")
		args = append(args, "%"+data.Description+"%")
	}

	// 构建完整 SQL
	sql := fmt.Sprintf("select %s from %s", positionRows, m.table)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM  %s", m.table)
	if len(conditions) > 0 {
		where := " where " + strings.Join(conditions, " OR ")
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

	err = m.QueryRowsNoCacheCtx(ctx, &position, sql, args...)

	switch {
	case err == nil:
		return position, total, nil
	case errors.Is(err, sqlc.ErrNotFound):
		// 搜索时数据为空不是错误，返回空列表
		return []*Position{}, total, nil
	default:
		return nil, 0, err
	}
}

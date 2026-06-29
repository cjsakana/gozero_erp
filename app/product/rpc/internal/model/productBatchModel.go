package model

import (
	"context"
	"erp/app/product/rpc/internal/types"
	"erp/common/xtime"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
)

var _ ProductBatchModel = (*customProductBatchModel)(nil)

type (
	// ProductBatchModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductBatchModel.
	ProductBatchModel interface {
		productBatchModel
		Search(ctx context.Context, data *types.SearchProductBatchParams) ([]*ProductBatch, int64, error)
	}

	customProductBatchModel struct {
		*defaultProductBatchModel
	}
)

// NewProductBatchModel returns a model for the database table.
func NewProductBatchModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ProductBatchModel {
	return &customProductBatchModel{
		defaultProductBatchModel: newProductBatchModel(conn, c, opts...),
	}
}

func (m *customProductBatchModel) Search(ctx context.Context, data *types.SearchProductBatchParams) ([]*ProductBatch, int64, error) {
	var productBatchs []*ProductBatch

	conditions := []string{}
	args := []any{}

	if data.ProductId != 0 {
		conditions = append(conditions, "product_id = ?")
		args = append(args, data.ProductId)
	}
	if data.BatchNo != "" {
		conditions = append(conditions, "batch_no = ?")
		args = append(args, data.BatchNo)
	}
	if !xtime.IsZeroTime(data.StartDate) {
		conditions = append(conditions, "production_date >= ?")
		args = append(args, data.StartDate)
	}

	if !xtime.IsZeroTime(data.EndDate) {
		conditions = append(conditions, "production_date <= ?")
		args = append(args, data.EndDate)
	}

	// 构建完整 SQL
	sql := fmt.Sprintf("select %s from %s", productBatchRows, m.table)
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
	err = m.QueryRowsNoCacheCtx(ctx, &productBatchs, sql, args...)

	switch {
	case err == nil:
		return productBatchs, total, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, 0, ErrNotFound
	default:
		return nil, 0, err
	}
}

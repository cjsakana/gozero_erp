package model

import (
	"context"
	"database/sql"
	"erp/app/product/rpc/internal/types"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
)

var _ ProductModel = (*customProductModel)(nil)

type (
	// ProductModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductModel.
	ProductModel interface {
		productModel
		XUpdate(ctx context.Context, newData *Product) error
		Search(ctx context.Context, data *types.SearchProductParams) ([]*Product, int64, error)
	}

	customProductModel struct {
		*defaultProductModel
	}
)

// NewProductModel returns a model for the database table.
func NewProductModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ProductModel {
	return &customProductModel{
		defaultProductModel: newProductModel(conn, c, opts...),
	}
}

func (m *customProductModel) XUpdate(ctx context.Context, newData *Product) error {
	var setClauses []string
	var args []interface{}

	if newData.ProductName != "" {
		setClauses = append(setClauses, "product_name = ?")
		args = append(args, newData.ProductName)
	}
	if newData.CategoryId != 0 {
		setClauses = append(setClauses, "category_id = ?")
		args = append(args, newData.CategoryId)
	}
	if newData.Specifications.String != "" {
		setClauses = append(setClauses, "specifications = ?")
		args = append(args, newData.Specifications.String)
	}
	if newData.Unit != "" {
		setClauses = append(setClauses, "unit = ?")
		args = append(args, newData.Unit)
	}
	if newData.PurchasePrice.Float64 != 0 {
		setClauses = append(setClauses, "purchase_price = ?")
		args = append(args, newData.PurchasePrice.Float64)
	}
	if newData.SellingPrice.Float64 != 0 {
		setClauses = append(setClauses, "selling_price = ?")
		args = append(args, newData.SellingPrice.Float64)
	}
	if newData.IsActive != 0 {
		setClauses = append(setClauses, "is_active = ?")
		args = append(args, newData.IsActive)
	}
	if newData.IsMaterial != 0 {
		setClauses = append(setClauses, "is_material = ?")
		args = append(args, newData.IsMaterial)
	}

	setClauses = append(setClauses, " updated_by  = ?")
	args = append(args, newData.UpdatedBy)

	erpProductProductIdKey := fmt.Sprintf("%s%v", cacheErpProductProductIdPrefix, newData.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, strings.Join(setClauses, ", "))
		args = append(args, newData.Id)
		return conn.ExecCtx(ctx, query, args...)
	}, erpProductProductIdKey)
	return err
}

func (m *customProductModel) Search(ctx context.Context, data *types.SearchProductParams) ([]*Product, int64, error) {
	var product []*Product

	conditions := []string{}
	args := []any{}

	if data.ProductNo != "" {
		conditions = append(conditions, "product_no = ?")
		args = append(args, data.ProductNo)
	}
	if data.ProductName != "" {
		conditions = append(conditions, "product_name like ?")
		args = append(args, "%"+data.ProductName+"%")
	}
	if data.CategoryId != 0 {
		conditions = append(conditions, "category_id = ?")
		args = append(args, data.CategoryId)
	}

	if data.IsActive != 0 {
		conditions = append(conditions, "is_active = ?")
		args = append(args, data.IsActive)
	}
	if data.IsMaterial != 0 {
		conditions = append(conditions, "is_material = ?")
		args = append(args, data.IsMaterial)
	}

	// 构建完整 SQL
	sql := fmt.Sprintf("select %s from %s", productRows, m.table)
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
	err = m.QueryRowsNoCacheCtx(ctx, &product, sql, args...)

	switch {
	case err == nil:
		return product, total, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, 0, ErrNotFound
	default:
		return nil, 0, err
	}
}

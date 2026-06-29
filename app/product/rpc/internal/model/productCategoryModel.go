package model

import (
	"context"
	"database/sql"
	types2 "erp/app/product/rpc/internal/types"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
)

var _ ProductCategoryModel = (*customProductCategoryModel)(nil)

type (
	// ProductCategoryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductCategoryModel.
	ProductCategoryModel interface {
		productCategoryModel
		XUpdate(ctx context.Context, newData *ProductCategory) error
		Search(ctx context.Context, data *types2.SearchProductCategoryParams) ([]*ProductCategory, int64, error)
		XDelete(ctx context.Context, categoryId int64) error
	}

	customProductCategoryModel struct {
		*defaultProductCategoryModel
	}
)

// NewProductCategoryModel returns a model for the database table.
func NewProductCategoryModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ProductCategoryModel {
	return &customProductCategoryModel{
		defaultProductCategoryModel: newProductCategoryModel(conn, c, opts...),
	}
}

func (m *customProductCategoryModel) XUpdate(ctx context.Context, newData *ProductCategory) error {
	var setClauses []string
	var args []interface{}

	if newData.Name != "" {
		setClauses = append(setClauses, "name = ?")
		args = append(args, newData.Name)
	}
	if newData.ParentId != 0 {
		setClauses = append(setClauses, "parent_id = ?")
		args = append(args, newData.ParentId)
	}

	erpProductProductCategoryCategoryIdKey := fmt.Sprintf("%s%v", cacheErpProductProductCategoryIdPrefix, newData.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, strings.Join(setClauses, ", "))
		args = append(args, newData.Id)
		return conn.ExecCtx(ctx, query, args...)
	}, erpProductProductCategoryCategoryIdKey)
	return err
}

func (m *customProductCategoryModel) Search(ctx context.Context, data *types2.SearchProductCategoryParams) ([]*ProductCategory, int64, error) {
	var productCategories []*ProductCategory

	conditions := []string{}
	args := []any{}

	if data.CategoryName != "" {
		conditions = append(conditions, "name like ?")
		args = append(args, "%"+data.CategoryName+"%")
	}
	if data.ParentId != 0 {
		conditions = append(conditions, "parent_id = ?")
		args = append(args, data.ParentId)
	}

	// 构建完整 SQL
	sql := fmt.Sprintf("select %s from %s", productCategoryRows, m.table)
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
	err = m.QueryRowsNoCacheCtx(ctx, &productCategories, sql, args...)

	switch {
	case err == nil:
		return productCategories, total, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, 0, ErrNotFound
	default:
		return nil, 0, err
	}
}

func (m *customProductCategoryModel) XDelete(ctx context.Context, categoryId int64) error {
	// 1. 检查是否存在子分类
	_, total, err := m.Search(ctx, &types2.SearchProductCategoryParams{
		SearchCom: types2.SearchCom{
			Limit: -1,
		},
		ParentId: categoryId,
	})

	if err != nil && !errors.Is(err, sqlc.ErrNotFound) {
		// 有错误且不是"未找到"的错误
		return fmt.Errorf("检查子分类失败: %v", err)
	}

	if total > 0 {
		// 有子分类（无论是否有ErrNotFound错误）
		return errors.New("存在子分类，无法删除")
	}

	// 2. 检查是否有商品使用该分类
	var productCount int
	err = m.QueryRowNoCacheCtx(ctx, &productCount, "SELECT COUNT(*) FROM product WHERE category_id = ?", categoryId)
	if err != nil && !errors.Is(err, sqlc.ErrNotFound) {
		return fmt.Errorf("检查商品数量失败: %v", err)
	}
	if productCount > 0 {
		return fmt.Errorf("有%d个商品使用该分类，无法删除", productCount)
	}

	// 3. 删除分类
	erpProductProductCategoryCategoryIdKey := fmt.Sprintf("%s%v", cacheErpProductProductCategoryIdPrefix, categoryId)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, categoryId)
	}, erpProductProductCategoryCategoryIdKey)

	return err
}

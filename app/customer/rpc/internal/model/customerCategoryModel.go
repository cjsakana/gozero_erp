package model

import (
	"context"
	"database/sql"
	"erp/app/customer/rpc/internal/types"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
)

var _ CustomerCategoryModel = (*customCustomerCategoryModel)(nil)

type (
	// CustomerCategoryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCustomerCategoryModel.
	CustomerCategoryModel interface {
		customerCategoryModel
		XUpdate(ctx context.Context, data *CustomerCategory) error
		Search(ctx context.Context, data *types.SearchCustomerCategory) ([]*CustomerCategory, int64, error)
		XDelete(ctx context.Context, id int64) error
	}

	customCustomerCategoryModel struct {
		*defaultCustomerCategoryModel
	}
)

// NewCustomerCategoryModel returns a model for the database table.
func NewCustomerCategoryModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CustomerCategoryModel {
	return &customCustomerCategoryModel{
		defaultCustomerCategoryModel: newCustomerCategoryModel(conn, c, opts...),
	}
}

func (m *customCustomerCategoryModel) XUpdate(ctx context.Context, data *CustomerCategory) error {
	var setClauses []string
	var args []interface{}

	if data.Name != "" {
		setClauses = append(setClauses, "name = ?")
		args = append(args, data.Name)
	}
	if data.CreditPolicy.String != "" {
		setClauses = append(setClauses, "credit_policy = ?")
		args = append(args, data.CreditPolicy.String)
	}
	// 补上 id
	args = append(args, data.Id)

	erpCustomerCustomerCategoryIdKey := fmt.Sprintf("%s%v", cacheErpCustomerCustomerCategoryIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, strings.Join(setClauses, ", "))
		return conn.ExecCtx(ctx, query, args...)
	}, erpCustomerCustomerCategoryIdKey)
	return err
}

func (m *customCustomerCategoryModel) Search(ctx context.Context, data *types.SearchCustomerCategory) ([]*CustomerCategory, int64, error) {
	var customerCategories []*CustomerCategory

	conditions := []string{}
	args := []any{}

	if data.Name != "" {
		conditions = append(conditions, "name like ?")
		args = append(args, "%"+data.Name+"%")
	}
	if data.CreditPolicy != "" {
		conditions = append(conditions, "credit_policy like ?")
		args = append(args, "%"+data.CreditPolicy+"%")
	}

	// 构建完整 SQL
	sql := fmt.Sprintf("select %s from %s", customerCategoryRows, m.table)
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
	err = m.QueryRowsNoCacheCtx(ctx, &customerCategories, sql, args...)

	switch {
	case err == nil:
		return customerCategories, total, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, 0, ErrNotFound
	default:
		return nil, 0, err
	}
}

func (m *customCustomerCategoryModel) XDelete(ctx context.Context, id int64) error {
	// 检查是否有客户使用该分类
	var customerCount int
	err := m.QueryRowNoCacheCtx(ctx, &customerCount, "SELECT COUNT(*) FROM customer WHERE category_id = ?", id)
	if err != nil {
		return fmt.Errorf("检查客户失败: %v", err)
	}
	if customerCount > 0 {
		return fmt.Errorf("有%d个客户使用该分类，无法删除", customerCount)
	}

	//删除分类
	erpCustomerCustomerCategoryIdKey := fmt.Sprintf("%s%v", cacheErpCustomerCustomerCategoryIdPrefix, id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, erpCustomerCustomerCategoryIdKey)

	return err
}

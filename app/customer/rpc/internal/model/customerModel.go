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

var _ CustomerModel = (*customCustomerModel)(nil)

type (
	// CustomerModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCustomerModel.
	CustomerModel interface {
		customerModel
		XUpdate(ctx context.Context, newData *Customer) error
		Search(ctx context.Context, data *types.SearchCustomer) ([]*Customer, int64, error)
	}

	customCustomerModel struct {
		*defaultCustomerModel
	}
)

// NewCustomerModel returns a model for the database table.
func NewCustomerModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CustomerModel {
	return &customCustomerModel{
		defaultCustomerModel: newCustomerModel(conn, c, opts...),
	}
}

func (m *customCustomerModel) XUpdate(ctx context.Context, newData *Customer) error {
	var setClauses []string
	var args []interface{}

	if newData.Name != "" {
		setClauses = append(setClauses, "name = ?")
		args = append(args, newData.Name)
	}
	if newData.CategoryId != 0 {
		setClauses = append(setClauses, "category_id = ?")
		args = append(args, newData.CategoryId)
	}
	if newData.Contact.String != "" {
		setClauses = append(setClauses, "contact = ?")
		args = append(args, newData.Contact.String)
	}
	if newData.Phone.String != "" {
		setClauses = append(setClauses, "phone = ?")
		args = append(args, newData.Phone.String)
	}
	if newData.Address.String != "" {
		setClauses = append(setClauses, "address = ?")
		args = append(args, newData.Address.String)
	}

	if newData.CreditLimit.Float64 != 0 {
		setClauses = append(setClauses, "credit_limit = ?")
		args = append(args, newData.CreditLimit.Float64)
	}
	if newData.UsedCredit.Float64 != 0 {
		setClauses = append(setClauses, "used_credit = ?")
		args = append(args, newData.UsedCredit.Float64)
	}
	if newData.PaymentTerms.String != "" {
		setClauses = append(setClauses, "payment_terms = ?")
		args = append(args, newData.PaymentTerms.String)
	}
	if newData.IsActive != 0 {
		setClauses = append(setClauses, "is_active = ?")
		args = append(args, newData.IsActive)
	}

	setClauses = append(setClauses, "updated_by = ?")
	args = append(args, newData.UpdatedBy)

	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	erpCustomerCustomerCodeKey := fmt.Sprintf("%s%v", cacheErpCustomerCustomerCodePrefix, data.Code)
	erpCustomerCustomerUsccKey := fmt.Sprintf("%s%v", cacheErpCustomerCustomerUsccPrefix, data.Code)
	erpCustomerCustomerIdKey := fmt.Sprintf("%s%v", cacheErpCustomerCustomerIdPrefix, data.Id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, strings.Join(setClauses, ", "))
		args = append(args, newData.Id)
		return conn.ExecCtx(ctx, query, args...)
	}, erpCustomerCustomerCodeKey, erpCustomerCustomerUsccKey, erpCustomerCustomerIdKey)
	return err
}

func (m *customCustomerModel) Search(ctx context.Context, data *types.SearchCustomer) ([]*Customer, int64, error) {
	var productCategories []*Customer

	conditions := []string{}
	args := []any{}

	if data.Code != "" {
		conditions = append(conditions, "code = ?")
		args = append(args, data.Code)
	}
	if data.USCC != "" {
		conditions = append(conditions, "uscc = ?")
		args = append(args, data.USCC)
	}
	if data.Name != "" {
		conditions = append(conditions, "name like ?")
		args = append(args, "%"+data.Name+"%")
	}
	if data.CategoryId != 0 {
		conditions = append(conditions, "category_id = ?")
		args = append(args, data.CategoryId)
	}
	if data.Contact != "" {
		conditions = append(conditions, "contact like ?")
		args = append(args, "%"+data.Contact+"%")
	}
	if data.Address != "" {
		conditions = append(conditions, "address like ?")
		args = append(args, "%"+data.Address+"%")
	}
	if data.PaymentTerms != "" {
		conditions = append(conditions, "payment_terms like ?")
		args = append(args, "%"+data.PaymentTerms+"%")
	}
	if data.IsActive != 0 {
		conditions = append(conditions, "is_active = ?")
		args = append(args, data.IsActive)
	}

	// 构建完整 SQL
	sql := fmt.Sprintf("select %s from %s", customerRows, m.table)
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

package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SalaryPaymentModel = (*customSalaryPaymentModel)(nil)

type (
	// SalaryPaymentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSalaryPaymentModel.
	SalaryPaymentModel interface {
		salaryPaymentModel
		Search(ctx context.Context, employeeId, departmentId, status int64, paymentMonth int64, page, limit int64) ([]*SalaryPayment, int64, error)
	}

	customSalaryPaymentModel struct {
		*defaultSalaryPaymentModel
		conn sqlx.SqlConn
	}
)

// NewSalaryPaymentModel returns a model for the database table.
func NewSalaryPaymentModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SalaryPaymentModel {
	return &customSalaryPaymentModel{
		defaultSalaryPaymentModel: newSalaryPaymentModel(conn, c, opts...),
		conn:                      conn,
	}
}

func (m *customSalaryPaymentModel) Search(ctx context.Context, employeeId, departmentId, status int64, paymentMonth int64, page, limit int64) ([]*SalaryPayment, int64, error) {
	var salaryPayments []*SalaryPayment
	conditions := []string{}
	args := []any{}

	if employeeId != 0 {
		conditions = append(conditions, "employee_id = ?")
		args = append(args, employeeId)
	}
	if departmentId != 0 {
		conditions = append(conditions, "department_id = ?")
		args = append(args, departmentId)
	}
	if paymentMonth != 0 {
		conditions = append(conditions, "payment_month = DATE(FROM_UNIXTIME(?))")
		args = append(args, paymentMonth)
	}
	if status != 0 {
		conditions = append(conditions, "status = ?")
		args = append(args, status)
	}

	base := fmt.Sprintf("select %s from %s", salaryPaymentRows, m.table)
	countBase := fmt.Sprintf("select count(*) from %s", m.table)
	if len(conditions) > 0 {
		where := " where " + strings.Join(conditions, " AND ")
		base += where
		countBase += where
	}

	var total int64
	if err := m.QueryRowNoCacheCtx(ctx, &total, countBase, args...); err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		base += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, (page-1)*limit)
	}

	if err := m.QueryRowsNoCacheCtx(ctx, &salaryPayments, base, args...); err != nil {
		return nil, 0, err
	}
	return salaryPayments, total, nil
}

package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PaymentRecordModel = (*customPaymentRecordModel)(nil)

type (
	// PaymentRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPaymentRecordModel.
	PaymentRecordModel interface {
		paymentRecordModel
		Search(ctx context.Context, paymentNo, paymentMethod string, supplierId, paymentType, status int64, page, limit int64) ([]*PaymentRecord, int64, error)
	}

	customPaymentRecordModel struct {
		*defaultPaymentRecordModel
		conn sqlx.SqlConn
	}
)

// NewPaymentRecordModel returns a model for the database table.
func NewPaymentRecordModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PaymentRecordModel {
	return &customPaymentRecordModel{
		defaultPaymentRecordModel: newPaymentRecordModel(conn, c, opts...),
		conn:                      conn,
	}
}

func (m *customPaymentRecordModel) Search(ctx context.Context, paymentNo, paymentMethod string, supplierId, paymentType, status int64, page, limit int64) ([]*PaymentRecord, int64, error) {
	var paymentRecords []*PaymentRecord
	conditions := []string{}
	args := []any{}

	if paymentNo != "" {
		conditions = append(conditions, "payment_no = ?")
		args = append(args, paymentNo)
	}
	if supplierId != 0 {
		conditions = append(conditions, "supplier_id = ?")
		args = append(args, supplierId)
	}
	if paymentType != 0 {
		conditions = append(conditions, "payment_type = ?")
		args = append(args, paymentType)
	}
	if paymentMethod != "" {
		conditions = append(conditions, "payment_method = ?")
		args = append(args, paymentMethod)
	}
	if status != 0 {
		conditions = append(conditions, "status = ?")
		args = append(args, status)
	}

	base := fmt.Sprintf("select %s from %s", paymentRecordRows, m.table)
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

	if err := m.QueryRowsNoCacheCtx(ctx, &paymentRecords, base, args...); err != nil {
		return nil, 0, err
	}
	return paymentRecords, total, nil
}

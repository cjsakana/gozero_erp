package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ReceiptRecordModel = (*customReceiptRecordModel)(nil)

type (
	// ReceiptRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customReceiptRecordModel.
	ReceiptRecordModel interface {
		receiptRecordModel
		Search(ctx context.Context, receiptNo, receiptMethod, operatorNo string, customerId, receiptType, status int64, page, limit int64) ([]*ReceiptRecord, int64, error)
	}

	customReceiptRecordModel struct {
		*defaultReceiptRecordModel
		conn sqlx.SqlConn
	}
)

// NewReceiptRecordModel returns a model for the database table.
func NewReceiptRecordModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ReceiptRecordModel {
	return &customReceiptRecordModel{
		defaultReceiptRecordModel: newReceiptRecordModel(conn, c, opts...),
		conn:                      conn,
	}
}

func (m *customReceiptRecordModel) Search(ctx context.Context, receiptNo, receiptMethod, operatorNo string, customerId, receiptType, status int64, page, limit int64) ([]*ReceiptRecord, int64, error) {
	var receiptRecords []*ReceiptRecord
	conditions := []string{}
	args := []any{}

	if receiptNo != "" {
		conditions = append(conditions, "receipt_no = ?")
		args = append(args, receiptNo)
	}
	if customerId != 0 {
		conditions = append(conditions, "customer_id = ?")
		args = append(args, customerId)
	}
	if receiptType != 0 {
		conditions = append(conditions, "receipt_type = ?")
		args = append(args, receiptType)
	}
	if receiptMethod != "" {
		conditions = append(conditions, "receipt_method = ?")
		args = append(args, receiptMethod)
	}
	if status != 0 {
		conditions = append(conditions, "status = ?")
		args = append(args, status)
	}
	if operatorNo != "" {
		conditions = append(conditions, "operator_no = ?")
		args = append(args, operatorNo)
	}

	base := fmt.Sprintf("select %s from %s", receiptRecordRows, m.table)
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

	if err := m.QueryRowsNoCacheCtx(ctx, &receiptRecords, base, args...); err != nil {
		return nil, 0, err
	}
	return receiptRecords, total, nil
}

package model

import (
	"context"
	"database/sql"
	"erp/app/purchase/rpc/internal/types"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PurchaseReceiptDetailModel = (*customPurchaseReceiptDetailModel)(nil)

type (
	// PurchaseReceiptDetailModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPurchaseReceiptDetailModel.
	PurchaseReceiptDetailModel interface {
		purchaseReceiptDetailModel
		ListByReceiptId(ctx context.Context, receiptId int64) ([]*PurchaseReceiptDetail, error)
		UpdateReceiptDetail(ctx context.Context, data *types.UpdateReceiptDetailParam) error
	}

	customPurchaseReceiptDetailModel struct {
		*defaultPurchaseReceiptDetailModel
	}
)

// NewPurchaseReceiptDetailModel returns a model for the database table.
func NewPurchaseReceiptDetailModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PurchaseReceiptDetailModel {
	return &customPurchaseReceiptDetailModel{
		defaultPurchaseReceiptDetailModel: newPurchaseReceiptDetailModel(conn, c, opts...),
	}
}

func (m *customPurchaseReceiptDetailModel) ListByReceiptId(ctx context.Context, receiptId int64) ([]*PurchaseReceiptDetail, error) {
	var list []*PurchaseReceiptDetail
	query := "select " + purchaseReceiptDetailRows + " from purchase_receipt_detail where receipt_id = ?"
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, receiptId); err != nil {
		return nil, err
	}
	return list, nil
}

// 更新采购入库明细（动态SQL拼接）
func (m *customPurchaseReceiptDetailModel) UpdateReceiptDetail(ctx context.Context, data *types.UpdateReceiptDetailParam) error {
	// 动态构建SQL
	setParts := []string{}
	args := []any{}

	if data.ProductId != nil {
		setParts = append(setParts, "product_id = ?")
		args = append(args, *data.ProductId)
	}
	if data.ProductName != nil {
		setParts = append(setParts, "product_name = ?")
		args = append(args, *data.ProductName)
	}
	if data.CategoryType != nil {
		setParts = append(setParts, "category_type = ?")
		args = append(args, *data.CategoryType)
	}
	if data.Quantity != nil {
		setParts = append(setParts, "quantity = ?")
		args = append(args, *data.Quantity)
	}
	if data.UnitPrice != nil {
		setParts = append(setParts, "unit_price = ?")
		args = append(args, *data.UnitPrice)
	}
	if data.Amount != nil {
		setParts = append(setParts, "amount = ?")
		args = append(args, *data.Amount)
	}
	if data.BatchId != nil {
		setParts = append(setParts, "batch_id = ?")
		args = append(args, *data.BatchId)
	}

	// 如果没有要更新的字段，直接返回
	if len(setParts) == 0 {
		return nil
	}

	args = append(args, data.Id) // WHERE条件的参数

	// 构建完整SQL
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", m.table, strings.Join(setParts, ", "))

	// 执行更新并清理缓存
	erpPurchasePurchaseReceiptDetailIdKey := fmt.Sprintf("%s%v", cacheErpPurchasePurchaseReceiptDetailIdPrefix, data.Id)
	
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, erpPurchasePurchaseReceiptDetailIdKey)
	
	return err
}

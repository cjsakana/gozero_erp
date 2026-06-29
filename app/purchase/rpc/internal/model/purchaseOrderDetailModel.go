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

var _ PurchaseOrderDetailModel = (*customPurchaseOrderDetailModel)(nil)

type (
	// PurchaseOrderDetailModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPurchaseOrderDetailModel.
	PurchaseOrderDetailModel interface {
		purchaseOrderDetailModel
		ListByOrderId(ctx context.Context, orderId int64) ([]*PurchaseOrderDetail, error)
		UpdateReceivedQty(ctx context.Context, orderId int64, productId int64, receivedQty float64) error
		UpdateOrderDetail(ctx context.Context, data *types.UpdateOrderDetailParam) error
	}

	customPurchaseOrderDetailModel struct {
		*defaultPurchaseOrderDetailModel
	}
)

// NewPurchaseOrderDetailModel returns a model for the database table.
func NewPurchaseOrderDetailModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PurchaseOrderDetailModel {
	return &customPurchaseOrderDetailModel{
		defaultPurchaseOrderDetailModel: newPurchaseOrderDetailModel(conn, c, opts...),
	}
}

func (m *customPurchaseOrderDetailModel) ListByOrderId(ctx context.Context, orderId int64) ([]*PurchaseOrderDetail, error) {
	var list []*PurchaseOrderDetail
	query := "select " + purchaseOrderDetailRows + " from purchase_order_detail where order_id = ?"
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, orderId); err != nil {
		return nil, err
	}
	return list, nil
}

func (m *customPurchaseOrderDetailModel) UpdateReceivedQty(ctx context.Context, orderId int64, productId int64, receivedQty float64) error {
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := "UPDATE purchase_order_detail SET received_qty = ? WHERE order_id = ? AND product_id = ?"
		return conn.ExecCtx(ctx, query, receivedQty, orderId, productId)
	})
	return err
}

// 更新采购订单明细（动态SQL拼接）
func (m *customPurchaseOrderDetailModel) UpdateOrderDetail(ctx context.Context, data *types.UpdateOrderDetailParam) error {
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
	if data.ReceivedQty != nil {
		setParts = append(setParts, "received_qty = ?")
		args = append(args, *data.ReceivedQty)
	}
	if data.Remark != nil {
		setParts = append(setParts, "remark = ?")
		args = append(args, *data.Remark)
	}

	// 如果没有要更新的字段，直接返回
	if len(setParts) == 0 {
		return nil
	}

	args = append(args, data.Id) // WHERE条件的参数

	// 构建完整SQL
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", m.table, strings.Join(setParts, ", "))

	// 执行更新并清理缓存
	erpPurchasePurchaseOrderDetailIdKey := fmt.Sprintf("%s%v", cacheErpPurchasePurchaseOrderDetailIdPrefix, data.Id)
	
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, erpPurchasePurchaseOrderDetailIdKey)
	
	return err
}

package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SalesOrderDetailModel = (*customSalesOrderDetailModel)(nil)

type (
	// SalesOrderDetailModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSalesOrderDetailModel.
	SalesOrderDetailModel interface {
		salesOrderDetailModel
		ListByOrderId(ctx context.Context, orderId int64) ([]*SalesOrderDetail, error)
		XUpdate(ctx context.Context, data *SalesOrderDetail) error
	}

	customSalesOrderDetailModel struct {
		*defaultSalesOrderDetailModel
	}
)

// NewSalesOrderDetailModel returns a model for the database table.
func NewSalesOrderDetailModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SalesOrderDetailModel {
	return &customSalesOrderDetailModel{
		defaultSalesOrderDetailModel: newSalesOrderDetailModel(conn, c, opts...),
	}
}

func (m *customSalesOrderDetailModel) ListByOrderId(ctx context.Context, orderId int64) ([]*SalesOrderDetail, error) {
	var list []*SalesOrderDetail
	query := "select " + salesOrderDetailRows + " from sales_order_detail where order_id = ?"
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, orderId); err != nil {
		return nil, err
	}
	return list, nil
}

// XUpdate 更新销售订单明细
func (m *customSalesOrderDetailModel) XUpdate(ctx context.Context, data *SalesOrderDetail) error {
	erpSalesSalesOrderDetailIdKey := fmt.Sprintf("%s%v", cacheErpSalesSalesOrderDetailIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := "update sales_order_detail set `order_id` = ?, `product_id` = ?, `product_name` = ?, `unit` = ?, `quantity` = ?, `unit_price` = ?, `amount` = ?, `delivered_qty` = ?, `remark` = ? where `id` = ?"
		return conn.ExecCtx(ctx, query, data.OrderId, data.ProductId, data.ProductName, data.Unit, data.Quantity, data.UnitPrice, data.Amount, data.DeliveredQty, data.Remark, data.Id)
	}, erpSalesSalesOrderDetailIdKey)
	return err
}

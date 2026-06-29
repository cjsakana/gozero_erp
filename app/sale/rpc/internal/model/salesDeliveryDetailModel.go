package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SalesDeliveryDetailModel = (*customSalesDeliveryDetailModel)(nil)

type (
	// SalesDeliveryDetailModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSalesDeliveryDetailModel.
	SalesDeliveryDetailModel interface {
		salesDeliveryDetailModel
		FindByDeliveryId(ctx context.Context, deliveryId int64) ([]*SalesDeliveryDetail, error)
		XUpdate(ctx context.Context, data *SalesDeliveryDetail) error
	}

	customSalesDeliveryDetailModel struct {
		*defaultSalesDeliveryDetailModel
	}
)

// NewSalesDeliveryDetailModel returns a model for the database table.
func NewSalesDeliveryDetailModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SalesDeliveryDetailModel {
	return &customSalesDeliveryDetailModel{
		defaultSalesDeliveryDetailModel: newSalesDeliveryDetailModel(conn, c, opts...),
	}
}

func (m *customSalesDeliveryDetailModel) FindByDeliveryId(ctx context.Context, deliveryId int64) ([]*SalesDeliveryDetail, error) {

	var detail []*SalesDeliveryDetail

	query := fmt.Sprintf("select %s from %s where `delivery_id` = ?", salesDeliveryDetailRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &detail, query, deliveryId)
	switch err {
	case nil:
		return detail, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// XUpdate 更新销售出库明细
func (m *customSalesDeliveryDetailModel) XUpdate(ctx context.Context, data *SalesDeliveryDetail) error {
	erpSalesSalesDeliveryDetailIdKey := fmt.Sprintf("%s%v", cacheErpSalesSalesDeliveryDetailIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set `delivery_id` = ?, `product_id` = ?, `product_name` = ?, `unit` = ?, `quantity` = ?, `unit_price` = ?, `amount` = ?, `batch_id` = ? where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, data.DeliveryId, data.ProductId, data.ProductName, data.Unit, data.Quantity, data.UnitPrice, data.Amount, data.BatchId, data.Id)
	}, erpSalesSalesDeliveryDetailIdKey)
	return err
}

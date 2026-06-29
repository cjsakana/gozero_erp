package model

import (
	"context"
	"erp/app/sale/rpc/internal/types"
	"erp/common/xtime"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SalesDeliveryModel = (*customSalesDeliveryModel)(nil)

type (
	// SalesDeliveryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSalesDeliveryModel.
	SalesDeliveryModel interface {
		salesDeliveryModel
		AddWithDetails(ctx context.Context, data *types.AddSalesDeliveryParam) error
		Outbound(ctx context.Context, data *types.OutboundParam) ([]string, error)
		Search(ctx context.Context, data *types.SearchDeliveryParams) ([]*SalesDelivery, int64, error)
	}

	customSalesDeliveryModel struct {
		*defaultSalesDeliveryModel
		conn sqlx.SqlConn
	}
)

// NewSalesDeliveryModel returns a model for the database table.
func NewSalesDeliveryModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SalesDeliveryModel {
	return &customSalesDeliveryModel{
		defaultSalesDeliveryModel: newSalesDeliveryModel(conn, c, opts...),
		conn:                      conn,
	}
}

func (m *customSalesDeliveryModel) AddWithDetails(ctx context.Context, data *types.AddSalesDeliveryParam) error {
	err := m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		// 1. 插入出库主表
		_, err := session.ExecCtx(ctx,
			`INSERT INTO sales_delivery (id, delivery_no, order_id, warehouse_id, delivery_date, total_quantity, total_amount, created_by)
			 VALUES (?, ?, ?, ?, FROM_UNIXTIME(?), ?, ?, ?)`,
			data.Id, data.DeliveryNo, data.OrderId, data.WarehouseId, data.DeliveryDate, data.TotalQuantity, data.TotalAmount, data.CreatedBy)
		if err != nil {
			fmt.Println("11111", err)
			return err
		}

		// 2. 插入明细表
		for _, item := range data.Details {
			if item.BatchId > 0 {
				_, err = session.ExecCtx(ctx,
					`INSERT INTO sales_delivery_detail (id,delivery_id, product_id, product_name, unit, quantity, unit_price, amount,batch_id)
				 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
					item.Id, item.DeliveryId, item.ProductId, item.ProductName, item.Unit, item.Quantity, item.UnitPrice, item.Amount, item.BatchId)
				if err != nil {
					fmt.Println("222222", err)

					return err
				}
			} else {
				_, err = session.ExecCtx(ctx,
					`INSERT INTO sales_delivery_detail (id,delivery_id, product_id, product_name, unit, quantity, unit_price, amount)
				 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
					item.Id, item.DeliveryId, item.ProductId, item.ProductName, item.Unit, item.Quantity, item.UnitPrice, item.Amount)
				if err != nil {
					fmt.Println("33333", err)

					return err
				}
			}
		}

		return nil
	})
	return err
}

func (m *customSalesDeliveryModel) Outbound(ctx context.Context, data *types.OutboundParam) ([]string, error) {
	// 缓存 key
	var keys []string

	// 查询 SalesDelivery
	erpSaleSalesDeliveryIdKey := fmt.Sprintf(types.CacheErpSaleSalesDeliveryIdPrefix, data.Id)
	var salesDelivery SalesDelivery
	_ = m.QueryRowCtx(ctx, &salesDelivery, erpSaleSalesDeliveryIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from sales_delivery where `id` = ? limit 1", salesDeliveryRows)
		return conn.QueryRowCtx(ctx, v, query, data.Id)
	})
	keys = append(keys, erpSaleSalesDeliveryIdKey)

	// 查询 SalesOrder
	erpSaleSalesOrderIdKey := fmt.Sprintf(types.CacheErpSaleSalesOrderIdPrefix, salesDelivery.OrderId.Int64)
	var salesOrder SalesOrder
	err := m.QueryRowCtx(ctx, &salesOrder, erpSaleSalesOrderIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from sales_order where `id` = ? and `status` !=4 limit 1", salesOrderRows)
		return conn.QueryRowCtx(ctx, v, query, salesDelivery.OrderId.Int64)
	})
	if errors.Is(err, sqlc.ErrNotFound) {
		return []string{}, fmt.Errorf("订单已取消")
	}
	keys = append(keys, erpSaleSalesOrderIdKey)

	// 其他错误会导致事务失败，所以我没处理了

	err = m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		// 1. 更新主表状态 + 操作人 + 出库时间
		_, err := session.ExecCtx(ctx, `
			UPDATE sales_delivery 
			SET status = 1, created_by = ?, delivery_date = ?
			WHERE id = ?`,
			data.CreatedBy, time.Now(), data.Id)
		if err != nil {
			return err
		}

		for _, item := range data.Items {
			//  更新该仓库交付明细的批次号，出库数量
			_, err := session.ExecCtx(ctx, `
				UPDATE sales_delivery_detail 
				SET batch_id = ?, quantity = ?
				WHERE id = ?`,
				item.BatchId, item.Quantity, item.Id)
			if err != nil {
				return err
			}
			keys = append(keys, fmt.Sprintf(types.CacheErpSaleSalesDeliveryDetailIdPrefix, item.Id))

			// 更新订单明细的出库数量
			_, err = session.ExecCtx(ctx, `UPDATE sales_order_detail SET delivered_qty = delivered_qty + ? 
                          WHERE order_id = ? and product_id = ?`,
				item.Quantity, salesOrder.Id, item.ProductId)
			if err != nil {
				return err
			}
		}

		// 标记到底多少商品完成了出库
		orderFinish := 0
		// 查询 SalesOrderDetail
		// 统计情况
		type Qty struct {
			Id           int64
			Quantity     float64
			DeliveredQty float64
		}
		var qty []*Qty
		query := fmt.Sprintf("select %s from sales_order_detail where `order_id` = ?", salesOrderDetailRows)
		err = session.QueryRowsCtx(ctx, &qty, query, salesDelivery.OrderId)
		if err != nil {
			return err
		}
		for _, item := range qty {
			if item.Quantity == item.DeliveredQty {
				orderFinish++
			}
			keys = append(keys, fmt.Sprintf(types.CacheErpSaleSalesOrderDetailIdPrefix, item.Id))
		}

		// 3.4 判断并更新订单状态
		// 状态： 1-已确认 2-部分发货 3-已完成 4-已取消
		status := 2
		if orderFinish == len(qty) {
			status = 3
		}
		session.ExecCtx(ctx, `UPDATE sales_order SET status = ? WHERE id = ?`, status, salesOrder.Id)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (m *customSalesDeliveryModel) Search(ctx context.Context, data *types.SearchDeliveryParams) ([]*SalesDelivery, int64, error) {
	var salesDeliveries []*SalesDelivery

	conditions := []string{}
	args := []any{}

	if data.DeliveryNo != "" {
		conditions = append(conditions, "delivery_no = ?")
		args = append(args, data.DeliveryNo)
	}
	if data.OrderId != 0 {
		conditions = append(conditions, "order_id = ?")
		args = append(args, data.OrderId)
	}
	if data.WarehouseId > 0 {
		conditions = append(conditions, "warehouse_id = ?")
		args = append(args, data.WarehouseId)
	}
	if !xtime.IsZeroTime(data.DeliveryDate) {
		conditions = append(conditions, "delivery_date >= ?")
		args = append(args, data.DeliveryDate)
	}
	if data.Status != 0 {
		conditions = append(conditions, "status = ?")
		args = append(args, data.Status)
	}

	// 构建完整 SQL
	sql := fmt.Sprintf("select %s from %s", salesDeliveryRows, m.table)
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
	err = m.QueryRowsNoCacheCtx(ctx, &salesDeliveries, sql, args...)

	switch {
	case err == nil:
		return salesDeliveries, total, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, 0, ErrNotFound
	default:
		return nil, 0, err
	}
}

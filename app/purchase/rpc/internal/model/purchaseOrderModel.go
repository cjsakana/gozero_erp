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

var _ PurchaseOrderModel = (*customPurchaseOrderModel)(nil)

type (
	// PurchaseOrderModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPurchaseOrderModel.
	PurchaseOrderModel interface {
		purchaseOrderModel
		CreateWithDetails(ctx context.Context, orderId int64, data *types.CreateOrderWithDetailsParam) error
		CreateFromRequisition(ctx context.Context, orderId int64, data *types.CreateOrderFromRequisitionParam) error
		CancelOrder(ctx context.Context, id int64) error
		Search(ctx context.Context, data *types.SearchOrderParams) ([]*PurchaseOrder, int64, error)
		UpdateContractURL(ctx context.Context, id int64, url string) error
		UpdateOrder(ctx context.Context, data *types.UpdateOrderParam) error
	}

	customPurchaseOrderModel struct {
		*defaultPurchaseOrderModel
		conn sqlx.SqlConn
	}
)

// NewPurchaseOrderModel returns a model for the database table.
func NewPurchaseOrderModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PurchaseOrderModel {
	return &customPurchaseOrderModel{
		defaultPurchaseOrderModel: newPurchaseOrderModel(conn, c, opts...),
		conn:                      conn,
	}
}

func (m *customPurchaseOrderModel) CreateWithDetails(ctx context.Context, orderId int64, data *types.CreateOrderWithDetailsParam) error {
	err := m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		_, err := session.ExecCtx(ctx,
			`INSERT INTO purchase_order (id, order_no, supplier_id, order_date, expected_date, total_amount, status, purchaser_id) 
			VALUES (?, ?, ?, FROM_UNIXTIME(?), FROM_UNIXTIME(?), ?, ?, ?)`,
			orderId, data.OrderNo, data.SupplierId, data.OrderDate, data.ExpectedDate, data.TotalAmount, data.Status, data.PurchaserId,
		)
		if err != nil {
			return err
		}

		// insert details
		if len(data.Details) > 0 {
			for _, d := range data.Details {
				_, err := session.ExecCtx(ctx,
					`INSERT INTO purchase_order_detail (id, order_id, product_id, product_name, category_type, quantity, unit_price, amount, received_qty, remark) 
					VALUES (?, ?, ?, ?, ?, ?, ?, ?, 0, ?)`,
					d.Id, orderId, d.ProductId, d.ProductName, d.CategoryType, d.Quantity, d.UnitPrice, d.Amount, d.Remark,
				)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}

func (m *customPurchaseOrderModel) CreateFromRequisition(ctx context.Context, orderId int64, data *types.CreateOrderFromRequisitionParam) error {
	// 这个方法需要先获取申请明细，然后创建订单
	// 实际实现需要访问 PurchaseRequisitionDetailModel，这里先提供一个基础实现
	// 在 logic 层处理会更合适
	err := m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		// 先插入订单主表，明细在 logic 层处理
		_, err := session.ExecCtx(ctx,
			`INSERT INTO purchase_order (id, order_no, supplier_id, order_date, expected_date, status, purchaser_id) 
			VALUES (?, ?, ?, FROM_UNIXTIME(?), FROM_UNIXTIME(?), 1, ?)`,
			orderId, data.OrderNo, data.SupplierId, data.OrderDate, data.ExpectedDate, data.PurchaserId,
		)
		if err != nil {
			return err
		}

		// 如果提供了明细覆盖，则使用覆盖的明细
		if len(data.Details) > 0 {
			var totalAmount float64
			for _, d := range data.Details {
				amount := d.Quantity * d.UnitPrice
				totalAmount += amount
				_, err := session.ExecCtx(ctx,
					`INSERT INTO purchase_order_detail (id, order_id, product_id, product_name, category_type, quantity, unit_price, amount, received_qty, remark) 
					VALUES (?, ?, ?, ?, ?, ?, ?, ?, 0, ?)`,
					d.Id, orderId, d.ProductId, d.ProductName, d.CategoryType, d.Quantity, d.UnitPrice, amount, d.Remark,
				)
				if err != nil {
					return err
				}
			}
			// 更新订单总金额
			_, err = session.ExecCtx(ctx, `UPDATE purchase_order SET total_amount = ? WHERE id = ?`, totalAmount, orderId)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (m *customPurchaseOrderModel) CancelOrder(ctx context.Context, id int64) error {
	_, err := m.conn.ExecCtx(ctx, "UPDATE purchase_order SET status = 5 WHERE id = ?", id)
	return err
}

func (m *customPurchaseOrderModel) Search(ctx context.Context, data *types.SearchOrderParams) ([]*PurchaseOrder, int64, error) {
	var orders []*PurchaseOrder
	conditions := []string{}
	args := []any{}

	if data.OrderNo != "" {
		conditions = append(conditions, "order_no = ?")
		args = append(args, data.OrderNo)
	}
	if data.SupplierId != 0 {
		conditions = append(conditions, "supplier_id = ?")
		args = append(args, data.SupplierId)
	}
	if data.Status != 0 {
		conditions = append(conditions, "status = ?")
		args = append(args, data.Status)
	}

	base := fmt.Sprintf("select %s from %s", purchaseOrderRows, m.table)
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

	base += " order by created_at desc"

	if data.Limit != -1 {
		base += fmt.Sprintf(" LIMIT %d OFFSET %d", data.Limit, (data.Page-1)*data.Limit)
	}

	if err := m.QueryRowsNoCacheCtx(ctx, &orders, base, args...); err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

func (m *customPurchaseOrderModel) UpdateContractURL(ctx context.Context, id int64, url string) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	erpPurchasePurchaseOrderIdKey := fmt.Sprintf("%s%v", cacheErpPurchasePurchaseOrderIdPrefix, data.Id)
	erpPurchasePurchaseOrderOrderNoKey := fmt.Sprintf("%s%v", cacheErpPurchasePurchaseOrderOrderNoPrefix, data.OrderNo)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("UPDATE %s SET contract_url = ? WHERE id = ?", m.table)
		return conn.ExecCtx(ctx, query, url, id)
	}, erpPurchasePurchaseOrderIdKey, erpPurchasePurchaseOrderOrderNoKey)
	return err
}

// 更新采购订单（动态SQL拼接）
func (m *customPurchaseOrderModel) UpdateOrder(ctx context.Context, data *types.UpdateOrderParam) error {
	// 获取原数据用于缓存清理
	original, err := m.FindOne(ctx, data.Id)
	if err != nil {
		return err
	}

	// 动态构建SQL
	setParts := []string{}
	args := []any{}

	if data.SupplierId != nil {
		setParts = append(setParts, "supplier_id = ?")
		args = append(args, *data.SupplierId)
	}
	if data.OrderDate != nil {
		setParts = append(setParts, "order_date = FROM_UNIXTIME(?)")
		args = append(args, *data.OrderDate)
	}
	if data.ExpectedDate != nil {
		setParts = append(setParts, "expected_date = FROM_UNIXTIME(?)")
		args = append(args, *data.ExpectedDate)
	}
	if data.TotalAmount != nil {
		setParts = append(setParts, "total_amount = ?")
		args = append(args, *data.TotalAmount)
	}
	if data.Status != nil {
		setParts = append(setParts, "status = ?")
		args = append(args, *data.Status)
	}
	if data.PurchaserId != nil {
		setParts = append(setParts, "purchaser_id = ?")
		args = append(args, *data.PurchaserId)
	}
	if data.ContractUrl != nil {
		setParts = append(setParts, "contract_url = ?")
		args = append(args, *data.ContractUrl)
	}

	// 如果没有要更新的字段，直接返回
	if len(setParts) == 0 {
		return nil
	}

	// 添加更新时间
	setParts = append(setParts, "updated_at = NOW()")
	args = append(args, data.Id) // WHERE条件的参数

	// 构建完整SQL
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", m.table, strings.Join(setParts, ", "))

	// 执行更新并清理缓存
	erpPurchasePurchaseOrderIdKey := fmt.Sprintf("%s%v", cacheErpPurchasePurchaseOrderIdPrefix, data.Id)
	erpPurchasePurchaseOrderOrderNoKey := fmt.Sprintf("%s%v", cacheErpPurchasePurchaseOrderOrderNoPrefix, original.OrderNo)
	
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, erpPurchasePurchaseOrderIdKey, erpPurchasePurchaseOrderOrderNoKey)
	
	return err
}

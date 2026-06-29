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

var _ PurchaseReceiptModel = (*customPurchaseReceiptModel)(nil)

type (
	// PurchaseReceiptModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPurchaseReceiptModel.
	PurchaseReceiptModel interface {
		purchaseReceiptModel
		CreateWithDetails(ctx context.Context, receiptId int64, data *types.CreateReceiptWithDetailsParam) error
		CreateFromOrder(ctx context.Context, receiptId int64, data *types.CreateReceiptFromOrderParam) error
		Search(ctx context.Context, data *types.SearchReceiptParams) ([]*PurchaseReceipt, int64, error)
		UpdateReceipt(ctx context.Context, data *types.UpdateReceiptParam) error
	}

	customPurchaseReceiptModel struct {
		*defaultPurchaseReceiptModel
		conn sqlx.SqlConn
	}
)

// NewPurchaseReceiptModel returns a model for the database table.
func NewPurchaseReceiptModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PurchaseReceiptModel {
	return &customPurchaseReceiptModel{
		defaultPurchaseReceiptModel: newPurchaseReceiptModel(conn, c, opts...),
		conn:                        conn,
	}
}

func (m *customPurchaseReceiptModel) CreateWithDetails(ctx context.Context, receiptId int64, data *types.CreateReceiptWithDetailsParam) error {
	err := m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		_, err := session.ExecCtx(ctx,
			`INSERT INTO purchase_receipt (id, receipt_no, order_id, warehouse_id, receipt_date, total_quantity, total_amount, status, created_by) 
			VALUES (?, ?, ?, ?, FROM_UNIXTIME(?), ?, ?, ?, ?)`,
			receiptId, data.ReceiptNo, sql.NullInt64{Int64: data.OrderId, Valid: data.OrderId > 0}, data.WarehouseId, data.ReceiptDate,
			data.TotalQuantity, data.TotalAmount, data.Status, data.CreatedBy,
		)
		if err != nil {
			return err
		}

		// insert details
		if len(data.Details) > 0 {
			for _, d := range data.Details {
				_, err := session.ExecCtx(ctx,
					`INSERT INTO purchase_receipt_detail (id, receipt_id, product_id, product_name, category_type, quantity, unit_price, amount, batch_id) 
					VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
					d.Id, receiptId, d.ProductId, d.ProductName, d.CategoryType, d.Quantity, d.UnitPrice, d.Amount, d.BatchId,
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

func (m *customPurchaseReceiptModel) CreateFromOrder(ctx context.Context, receiptId int64, data *types.CreateReceiptFromOrderParam) error {
	// 这个方法在 logic 层处理更合适，因为需要从订单明细获取数据
	// 这里提供一个基础实现框架
	err := m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		// 先插入入库单主表，明细在 logic 层处理
		var totalQuantity, totalAmount float64
		if len(data.Details) > 0 {
			for _, d := range data.Details {
				totalQuantity += d.Quantity
				totalAmount += d.Amount
			}
		}
		_, err := session.ExecCtx(ctx,
			`INSERT INTO purchase_receipt (id, receipt_no, order_id, warehouse_id, receipt_date, total_quantity, total_amount, status, created_by) 
			VALUES (?, ?, ?, ?, FROM_UNIXTIME(?), ?, ?, 1, ?)`,
			receiptId, data.ReceiptNo, sql.NullInt64{Int64: data.OrderId, Valid: data.OrderId > 0}, data.WarehouseId, data.ReceiptDate,
			totalQuantity, totalAmount, data.CreatedBy,
		)
		if err != nil {
			return err
		}

		// 如果提供了明细覆盖，则使用覆盖的明细
		if len(data.Details) > 0 {
			for _, d := range data.Details {
				_, err := session.ExecCtx(ctx,
					`INSERT INTO purchase_receipt_detail (id, receipt_id, product_id, product_name, category_type, quantity, unit_price, amount, batch_id) 
					VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
					d.Id, receiptId, d.ProductId, d.ProductName, d.CategoryType, d.Quantity, d.UnitPrice, d.Amount, d.BatchId,
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

func (m *customPurchaseReceiptModel) Search(ctx context.Context, data *types.SearchReceiptParams) ([]*PurchaseReceipt, int64, error) {
	var receipts []*PurchaseReceipt
	conditions := []string{}
	args := []any{}

	if data.ReceiptNo != "" {
		conditions = append(conditions, "receipt_no = ?")
		args = append(args, data.ReceiptNo)
	}
	if data.OrderId != 0 {
		conditions = append(conditions, "order_id = ?")
		args = append(args, data.OrderId)
	}
	if data.WarehouseId != 0 {
		conditions = append(conditions, "warehouse_id = ?")
		args = append(args, data.WarehouseId)
	}

	base := fmt.Sprintf("select %s from %s", purchaseReceiptRows, m.table)
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

	if err := m.QueryRowsNoCacheCtx(ctx, &receipts, base, args...); err != nil {
		return nil, 0, err
	}
	return receipts, total, nil
}

// 更新采购入库单（动态SQL拼接）
func (m *customPurchaseReceiptModel) UpdateReceipt(ctx context.Context, data *types.UpdateReceiptParam) error {
	// 获取原数据用于缓存清理
	original, err := m.FindOne(ctx, data.Id)
	if err != nil {
		return err
	}

	// 动态构建SQL
	setParts := []string{}
	args := []any{}

	if data.OrderId != nil {
		setParts = append(setParts, "order_id = ?")
		args = append(args, *data.OrderId)
	}
	if data.WarehouseId != nil {
		setParts = append(setParts, "warehouse_id = ?")
		args = append(args, *data.WarehouseId)
	}
	if data.ReceiptDate != nil {
		setParts = append(setParts, "receipt_date = FROM_UNIXTIME(?)")
		args = append(args, *data.ReceiptDate)
	}
	if data.TotalQuantity != nil {
		setParts = append(setParts, "total_quantity = ?")
		args = append(args, *data.TotalQuantity)
	}
	if data.TotalAmount != nil {
		setParts = append(setParts, "total_amount = ?")
		args = append(args, *data.TotalAmount)
	}
	if data.Status != nil {
		setParts = append(setParts, "status = ?")
		args = append(args, *data.Status)
	}
	if data.CreatedBy != nil {
		setParts = append(setParts, "created_by = ?")
		args = append(args, *data.CreatedBy)
	}

	// 如果没有要更新的字段，直接返回
	if len(setParts) == 0 {
		return nil
	}

	args = append(args, data.Id) // WHERE条件的参数

	// 构建完整SQL
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", m.table, strings.Join(setParts, ", "))

	// 执行更新并清理缓存
	erpPurchasePurchaseReceiptIdKey := fmt.Sprintf("%s%v", cacheErpPurchasePurchaseReceiptIdPrefix, data.Id)
	erpPurchasePurchaseReceiptReceiptNoKey := fmt.Sprintf("%s%v", cacheErpPurchasePurchaseReceiptReceiptNoPrefix, original.ReceiptNo)
	
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, erpPurchasePurchaseReceiptIdKey, erpPurchasePurchaseReceiptReceiptNoKey)
	
	return err
}

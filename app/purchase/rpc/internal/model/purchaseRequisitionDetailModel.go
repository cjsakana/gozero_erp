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

var _ PurchaseRequisitionDetailModel = (*customPurchaseRequisitionDetailModel)(nil)

type (
	// PurchaseRequisitionDetailModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPurchaseRequisitionDetailModel.
	PurchaseRequisitionDetailModel interface {
		purchaseRequisitionDetailModel
		ListByRequisitionId(ctx context.Context, requisitionId int64) ([]*PurchaseRequisitionDetail, error)
		UpdateRequisitionDetail(ctx context.Context, data *types.UpdateRequisitionDetailParam) error
	}

	customPurchaseRequisitionDetailModel struct {
		*defaultPurchaseRequisitionDetailModel
	}
)

// NewPurchaseRequisitionDetailModel returns a model for the database table.
func NewPurchaseRequisitionDetailModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PurchaseRequisitionDetailModel {
	return &customPurchaseRequisitionDetailModel{
		defaultPurchaseRequisitionDetailModel: newPurchaseRequisitionDetailModel(conn, c, opts...),
	}
}

func (m *customPurchaseRequisitionDetailModel) ListByRequisitionId(ctx context.Context, requisitionId int64) ([]*PurchaseRequisitionDetail, error) {
	var list []*PurchaseRequisitionDetail
	query := "select " + purchaseRequisitionDetailRows + " from purchase_requisition_detail where requisition_id = ?"
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, requisitionId); err != nil {
		return nil, err
	}
	return list, nil
}

// 更新采购申请明细（动态SQL拼接）
func (m *customPurchaseRequisitionDetailModel) UpdateRequisitionDetail(ctx context.Context, data *types.UpdateRequisitionDetailParam) error {

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
	erpPurchasePurchaseRequisitionDetailIdKey := fmt.Sprintf("%s%v", cacheErpPurchasePurchaseRequisitionDetailIdPrefix, data.Id)
	
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, erpPurchasePurchaseRequisitionDetailIdKey)
	
	return err
}

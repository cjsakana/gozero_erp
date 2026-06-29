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

var _ PurchaseRequisitionModel = (*customPurchaseRequisitionModel)(nil)

type (
	// PurchaseRequisitionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPurchaseRequisitionModel.
	PurchaseRequisitionModel interface {
		purchaseRequisitionModel
		CreateWithDetails(ctx context.Context, requisitionId int64, data *types.CreateRequisitionWithDetailsParam) error
		Approve(ctx context.Context, data *types.ApproveRequisitionParam) error
		Search(ctx context.Context, data *types.SearchRequisitionParams) ([]*PurchaseRequisition, int64, error)
		UpdateRequisition(ctx context.Context, data *types.UpdateRequisitionParam) error
	}

	customPurchaseRequisitionModel struct {
		*defaultPurchaseRequisitionModel
		conn sqlx.SqlConn
	}
)

// NewPurchaseRequisitionModel returns a model for the database table.
func NewPurchaseRequisitionModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PurchaseRequisitionModel {
	return &customPurchaseRequisitionModel{
		defaultPurchaseRequisitionModel: newPurchaseRequisitionModel(conn, c, opts...),
		conn:                            conn,
	}
}

func (m *customPurchaseRequisitionModel) CreateWithDetails(ctx context.Context, requisitionId int64, data *types.CreateRequisitionWithDetailsParam) error {
	err := m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		_, err := session.ExecCtx(ctx,
			`INSERT INTO purchase_requisition (id, requisition_no, department_id, applicant_id, approver_id, request_date, total_amount, status) 
			VALUES (?, ?, ?, ?, ?, FROM_UNIXTIME(?), ?, ?)`,
			requisitionId, data.RequisitionNo, data.DepartmentId, data.ApplicantId, data.ApproverId, data.RequestDate, data.TotalAmount, data.Status,
		)
		if err != nil {
			return err
		}

		// insert details
		if len(data.Details) > 0 {
			for _, d := range data.Details {
				_, err := session.ExecCtx(ctx,
					`INSERT INTO purchase_requisition_detail (id, requisition_id, product_id, product_name, category_type, quantity, unit_price, amount, remark) 
					VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
					d.Id, requisitionId, d.ProductId, d.ProductName, d.CategoryType, d.Quantity, d.UnitPrice, d.Amount, d.Remark,
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

func (m *customPurchaseRequisitionModel) Approve(ctx context.Context, newData *types.ApproveRequisitionParam) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	erpPurchasePurchaseRequisitionIdKey := fmt.Sprintf("%s%v", cacheErpPurchasePurchaseRequisitionIdPrefix, data.Id)
	erpPurchasePurchaseRequisitionRequisitionNoKey := fmt.Sprintf("%s%v", cacheErpPurchasePurchaseRequisitionRequisitionNoPrefix, data.RequisitionNo)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("UPDATE %s SET status = ?, approve_time = FROM_UNIXTIME(?), approve_remark = ? WHERE id = ?", m.table)
		return conn.ExecCtx(ctx, query, newData.TargetStatus, newData.ApproveTime, newData.ApproveRemark, newData.Id)
	}, erpPurchasePurchaseRequisitionIdKey, erpPurchasePurchaseRequisitionRequisitionNoKey)
	return err
}

func (m *customPurchaseRequisitionModel) Search(ctx context.Context, data *types.SearchRequisitionParams) ([]*PurchaseRequisition, int64, error) {
	var requisitions []*PurchaseRequisition
	conditions := []string{}
	args := []any{}

	if data.RequisitionNo != "" {
		conditions = append(conditions, "requisition_no = ?")
		args = append(args, data.RequisitionNo)
	}
	if data.DepartmentId != 0 {
		conditions = append(conditions, "department_id = ?")
		args = append(args, data.DepartmentId)
	}
	if data.ApplicantId != 0 {
		conditions = append(conditions, "applicant_id = ?")
		args = append(args, data.ApplicantId)
	}
	if data.ApproverId != 0 {
		conditions = append(conditions, "approver_id = ?")
		args = append(args, data.ApproverId)
	}
	if data.Status != 0 {
		conditions = append(conditions, "status = ?")
		args = append(args, data.Status)
	}

	base := fmt.Sprintf("select %s from %s", purchaseRequisitionRows, m.table)
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

	base += " order by created_at desc "

	if data.Limit != -1 {
		base += fmt.Sprintf(" LIMIT %d OFFSET %d", data.Limit, (data.Page-1)*data.Limit)
	}

	if err := m.QueryRowsNoCacheCtx(ctx, &requisitions, base, args...); err != nil {
		return nil, 0, err
	}
	return requisitions, total, nil
}

// 更新采购申请（动态SQL拼接）
func (m *customPurchaseRequisitionModel) UpdateRequisition(ctx context.Context, data *types.UpdateRequisitionParam) error {
	// 获取原数据用于缓存清理
	original, err := m.FindOne(ctx, data.Id)
	if err != nil {
		return err
	}

	// 动态构建SQL
	setParts := []string{}
	args := []any{}

	if data.DepartmentId != nil {
		setParts = append(setParts, "department_id = ?")
		args = append(args, *data.DepartmentId)
	}
	if data.ApplicantId != nil {
		setParts = append(setParts, "applicant_id = ?")
		args = append(args, *data.ApplicantId)
	}
	if data.RequestDate != nil {
		setParts = append(setParts, "request_date = FROM_UNIXTIME(?)")
		args = append(args, *data.RequestDate)
	}
	if data.TotalAmount != nil {
		setParts = append(setParts, "total_amount = ?")
		args = append(args, *data.TotalAmount)
	}
	if data.Status != nil {
		setParts = append(setParts, "status = ?")
		args = append(args, *data.Status)
	}
	if data.ApproverId != nil {
		setParts = append(setParts, "approver_id = ?")
		args = append(args, *data.ApproverId)
	}
	if data.ApproveTime != nil {
		setParts = append(setParts, "approve_time = FROM_UNIXTIME(?)")
		args = append(args, *data.ApproveTime)
	}
	if data.ApproveRemark != nil {
		setParts = append(setParts, "approve_remark = ?")
		args = append(args, *data.ApproveRemark)
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
	erpPurchasePurchaseRequisitionIdKey := fmt.Sprintf("%s%v", cacheErpPurchasePurchaseRequisitionIdPrefix, data.Id)
	erpPurchasePurchaseRequisitionRequisitionNoKey := fmt.Sprintf("%s%v", cacheErpPurchasePurchaseRequisitionRequisitionNoPrefix, original.RequisitionNo)
	
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, erpPurchasePurchaseRequisitionIdKey, erpPurchasePurchaseRequisitionRequisitionNoKey)
	
	return err
}

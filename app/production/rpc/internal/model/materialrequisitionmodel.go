package model

import (
	"context"
	"erp/app/production/rpc/internal/types"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MaterialRequisitionModel = (*customMaterialRequisitionModel)(nil)

type (
	// MaterialRequisitionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMaterialRequisitionModel.
	MaterialRequisitionModel interface {
		materialRequisitionModel
		GetRequisitionList(ctx context.Context, params *types.GetRequisitionListParams) ([]*MaterialRequisition, int64, error)
		CreateWithDetails(ctx context.Context, requisition *MaterialRequisition, items []*MaterialRequisitionItem) error
		UpdateWithDetails(ctx context.Context, requisition *MaterialRequisition, items []*MaterialRequisitionItem) error
	}

	customMaterialRequisitionModel struct {
		*defaultMaterialRequisitionModel
		conn sqlx.SqlConn
	}
)

// NewMaterialRequisitionModel returns a model for the database table.
func NewMaterialRequisitionModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) MaterialRequisitionModel {
	return &customMaterialRequisitionModel{
		defaultMaterialRequisitionModel: newMaterialRequisitionModel(conn, c, opts...),
		conn:                            conn,
	}
}

func (m *customMaterialRequisitionModel) CreateWithDetails(ctx context.Context, requisition *MaterialRequisition, items []*MaterialRequisitionItem) error {
	err := m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		query1 := fmt.Sprintf("insert into material_requisition (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", materialRequisitionRowsExpectAutoSet)
		_, err := session.ExecCtx(ctx, query1,
			requisition.Id, requisition.RequisitionNo, requisition.WorkOrderId, requisition.WorkOrderNo, requisition.WarehouseId, requisition.RequisitionDate,
			requisition.Status, requisition.CreatedBy, requisition.UpdatedBy,
		)
		if err != nil {
			return err
		}

		query2 := fmt.Sprintf("insert into material_requisition_item (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", materialRequisitionItemRowsExpectAutoSet)
		for _, d := range items {
			_, err := session.ExecCtx(ctx, query2,
				d.Id, d.RequisitionId, d.MaterialId, d.MaterialName, d.PlanQuantity, d.ActualQuantity, d.Unit, d.BatchNo, d.Remark,
			)
			if err != nil {
				return err
			}
		}

		return nil
	})
	return err
}

func (m *customMaterialRequisitionModel) GetRequisitionList(ctx context.Context, params *types.GetRequisitionListParams) ([]*MaterialRequisition, int64, error) {
	// 构造查询条件
	var conditions []string
	var args []interface{}

	if params.WorkOrderId != 0 {
		conditions = append(conditions, "work_order_id = ?")
		args = append(args, params.WorkOrderId)
	}
	if params.WarehouseId != 0 {
		conditions = append(conditions, "warehouse_id = ?")
		args = append(args, params.WarehouseId)
	}
	if params.Status != 0 {
		conditions = append(conditions, "status = ?")
		args = append(args, params.Status)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			whereClause += " AND " + conditions[i]
		}
	}

	// 查询总数
	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", m.table, whereClause)
	err := m.QueryRowNoCacheCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	offset := (params.Page - 1) * params.PageSize
	listQuery := fmt.Sprintf("SELECT * FROM %s %s ORDER BY created_at DESC LIMIT ? OFFSET ?", m.table, whereClause)
	args = append(args, params.PageSize, offset)

	var requisitions []*MaterialRequisition
	err = m.QueryRowsNoCacheCtx(ctx, &requisitions, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	return requisitions, total, nil
}

func (m *customMaterialRequisitionModel) UpdateWithDetails(ctx context.Context, requisition *MaterialRequisition, items []*MaterialRequisitionItem) error {
	err := m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		// 更新领料单主表
		query1 := fmt.Sprintf("update %s set `status` = ?, `updated_by` = ? where `id` = ?", m.table)
		_, err := session.ExecCtx(ctx, query1, requisition.Status, requisition.UpdatedBy, requisition.Id)
		if err != nil {
			return err
		}

		// 删除旧的明细
		query2 := "delete from material_requisition_item where `requisition_id` = ?"
		_, err = session.ExecCtx(ctx, query2, requisition.Id)
		if err != nil {
			return err
		}

		// 插入新的明细
		if len(items) > 0 {
			query3 := fmt.Sprintf("insert into material_requisition_item (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", materialRequisitionItemRowsExpectAutoSet)
			for _, item := range items {
				_, err := session.ExecCtx(ctx, query3,
					item.Id, item.RequisitionId, item.MaterialId, item.MaterialName, item.PlanQuantity, item.ActualQuantity, item.Unit, item.BatchNo, item.Remark,
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

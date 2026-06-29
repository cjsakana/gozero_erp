package model

import (
	"context"
	"erp/app/production/rpc/internal/types"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ WorkOrderModel = (*customWorkOrderModel)(nil)

type (
	// WorkOrderModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWorkOrderModel.
	WorkOrderModel interface {
		workOrderModel
		GetWorkOrderList(ctx context.Context, params *types.GetWorkOrderListParams) ([]*WorkOrder, int64, error)
	}

	customWorkOrderModel struct {
		*defaultWorkOrderModel
		conn sqlx.SqlConn
	}
)

// NewWorkOrderModel returns a model for the database table.
func NewWorkOrderModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) WorkOrderModel {
	return &customWorkOrderModel{
		defaultWorkOrderModel: newWorkOrderModel(conn, c, opts...),
		conn:                  conn,
	}
}

func (m *customWorkOrderModel) GetWorkOrderList(ctx context.Context, params *types.GetWorkOrderListParams) ([]*WorkOrder, int64, error) {
	// 构造查询条件
	var conditions []string
	var args []interface{}

	if params.ProductId != 0 {
		conditions = append(conditions, "product_id = ?")
		args = append(args, params.ProductId)
	}
	if params.Status != 0 {
		conditions = append(conditions, "status = ?")
		args = append(args, params.Status)
	}
	if params.Priority != 0 {
		conditions = append(conditions, "priority = ?")
		args = append(args, params.Priority)
	}
	if !params.StartDate.IsZero() && !params.EndDate.IsZero() {
		conditions = append(conditions, "plan_start_date >= ? AND plan_end_date <= ?")
		args = append(args, params.StartDate, params.EndDate)
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

	var workOrders []*WorkOrder
	err = m.QueryRowsNoCacheCtx(ctx, &workOrders, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	return workOrders, total, nil
}

package model

import (
	"context"
	"database/sql"
	"erp/app/hr/rpc/internal/types"
	"erp/common/xtime"
	"errors"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ResignedApplicationModel = (*customResignedApplicationModel)(nil)

type (
	// ResignedApplicationModel 离职申请模型接口
	ResignedApplicationModel interface {
		resignedApplicationModel
		XUpdate(ctx context.Context, newData *ResignedApplication) error
		Search(ctx context.Context, data *types.SearchResignedApplicationParams) ([]*ResignedApplication, int64, error)
	}

	customResignedApplicationModel struct {
		*defaultResignedApplicationModel
	}
)

// NewResignedApplicationModel 返回离职申请模型
func NewResignedApplicationModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ResignedApplicationModel {
	return &customResignedApplicationModel{
		defaultResignedApplicationModel: newResignedApplicationModel(conn, c, opts...),
	}
}

// XUpdate 动态更新离职申请（只更新非零值字段）
func (m *customResignedApplicationModel) XUpdate(ctx context.Context, newData *ResignedApplication) error {
	var setClauses []string
	var args []interface{}

	// 离职原因
	if newData.Reason != "" {
		setClauses = append(setClauses, "reason = ?")
		args = append(args, newData.Reason)
	}
	// 离职日期
	if !xtime.IsZeroTime(newData.LeaveDate) {
		setClauses = append(setClauses, "leave_date = ?")
		args = append(args, newData.LeaveDate)
	}
	// 证明材料
	if newData.Evidence.String != "" {
		setClauses = append(setClauses, "evidence = ?")
		args = append(args, newData.Evidence.String)
	}
	// 状态
	if newData.Status != 0 {
		setClauses = append(setClauses, "status = ?")
		args = append(args, newData.Status)
	}
	// 审批人ID
	if newData.ApproverId.Int64 != 0 {
		setClauses = append(setClauses, "approver_id = ?")
		args = append(args, newData.ApproverId.Int64)
	}
	// 审批时间
	if !xtime.IsZeroTime(newData.ApproveTime.Time) {
		setClauses = append(setClauses, "approve_time = ?")
		args = append(args, newData.ApproveTime.Time)
	}
	// 审批意见
	if newData.ApproveRemark.String != "" {
		setClauses = append(setClauses, "approve_remark = ?")
		args = append(args, newData.ApproveRemark.String)
	}

	erpHrResignedApplicationIdKey := fmt.Sprintf("%s%v", cacheErpHrResignedApplicationIdPrefix, newData.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, strings.Join(setClauses, ", "))
		args = append(args, newData.Id)
		return conn.ExecCtx(ctx, query, args...)
	}, erpHrResignedApplicationIdKey)
	return err
}

// Search 搜索离职申请
func (m *customResignedApplicationModel) Search(ctx context.Context, data *types.SearchResignedApplicationParams) ([]*ResignedApplication, int64, error) {
	var applications []*ResignedApplication

	conditions := []string{}
	args := []any{}

	// 申请人ID
	if data.EmployeeId != 0 {
		conditions = append(conditions, "employee_id = ?")
		args = append(args, data.EmployeeId)
	}
	// 离职日期范围
	if !xtime.IsZeroTime(data.StartLeaveDate) {
		conditions = append(conditions, "leave_date >= ?")
		args = append(args, data.StartLeaveDate)
	}
	if !xtime.IsZeroTime(data.EndLeaveDate) {
		conditions = append(conditions, "leave_date <= ?")
		args = append(args, data.EndLeaveDate)
	}
	// 状态
	if data.Status != 0 {
		conditions = append(conditions, "status = ?")
		args = append(args, data.Status)
	}
	// 审批人ID
	if data.ApproverId != 0 {
		conditions = append(conditions, "approver_id = ?")
		args = append(args, data.ApproverId)
	}

	// 构建完整 SQL
	sql := fmt.Sprintf("select %s from %s", resignedApplicationRows, m.table)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM  %s", m.table)
	if len(conditions) > 0 {
		where := " where " + strings.Join(conditions, " AND ")
		sql += where
		countQuery += where
	}

	sql += " Order by created_at desc "

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
	err = m.QueryRowsNoCacheCtx(ctx, &applications, sql, args...)

	switch {
	case err == nil:
		return applications, total, nil
	case errors.Is(err, sqlc.ErrNotFound):
		// 搜索时数据为空不是错误，返回空列表
		return []*ResignedApplication{}, total, nil
	default:
		return nil, 0, err
	}
}

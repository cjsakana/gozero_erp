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

var _ LeaveApplicationModel = (*customLeaveApplicationModel)(nil)

type (
	// LeaveApplicationModel 请假申请模型接口
	LeaveApplicationModel interface {
		leaveApplicationModel
		XUpdate(ctx context.Context, newData *LeaveApplication) error
		Search(ctx context.Context, data *types.SearchLeaveApplicationParams) ([]*LeaveApplication, int64, error)
	}

	customLeaveApplicationModel struct {
		*defaultLeaveApplicationModel
	}
)

// NewLeaveApplicationModel 返回请假申请模型
func NewLeaveApplicationModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) LeaveApplicationModel {
	return &customLeaveApplicationModel{
		defaultLeaveApplicationModel: newLeaveApplicationModel(conn, c, opts...),
	}
}

// XUpdate 动态更新请假申请（只更新非零值字段）
func (m *customLeaveApplicationModel) XUpdate(ctx context.Context, newData *LeaveApplication) error {
	var setClauses []string
	var args []interface{}

	// 请假类型
	if newData.Type != 0 {
		setClauses = append(setClauses, "type = ?")
		args = append(args, newData.Type)
	}
	// 开始时间
	if !xtime.IsZeroTime(newData.StartTime) {
		setClauses = append(setClauses, "start_time = ?")
		args = append(args, newData.StartTime)
	}
	// 结束时间
	if !xtime.IsZeroTime(newData.EndTime) {
		setClauses = append(setClauses, "end_time = ?")
		args = append(args, newData.EndTime)
	}
	// 时长
	if newData.Duration != 0 {
		setClauses = append(setClauses, "duration = ?")
		args = append(args, newData.Duration)
	}
	// 请假原因
	if newData.Reason != "" {
		setClauses = append(setClauses, "reason = ?")
		args = append(args, newData.Reason)
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

	erpHrLeaveApplicationIdKey := fmt.Sprintf("%s%v", cacheErpHrLeaveApplicationIdPrefix, newData.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, strings.Join(setClauses, ", "))
		args = append(args, newData.Id)
		return conn.ExecCtx(ctx, query, args...)
	}, erpHrLeaveApplicationIdKey)
	return err
}

// Search 搜索请假申请
func (m *customLeaveApplicationModel) Search(ctx context.Context, data *types.SearchLeaveApplicationParams) ([]*LeaveApplication, int64, error) {
	var leaveApplication []*LeaveApplication

	conditions := []string{}
	args := []any{}

	// 按员工ID查询
	if data.EmployeeId != 0 {
		conditions = append(conditions, "employee_id = ?")
		args = append(args, data.EmployeeId)
	}
	// 请假类型
	if data.Type != 0 {
		conditions = append(conditions, "type = ?")
		args = append(args, data.Type)
	}
	// 开始时间范围
	if !xtime.IsZeroTime(data.StartTime) {
		conditions = append(conditions, "start_time >= ?")
		args = append(args, data.StartTime)
	}
	// 结束时间范围
	if !xtime.IsZeroTime(data.EndTime) {
		conditions = append(conditions, "end_time <= ?")
		args = append(args, data.EndTime)
	}
	// 请假原因模糊查询
	if data.Reason != "" {
		conditions = append(conditions, "reason like ?")
		args = append(args, "%"+data.Reason+"%")
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
	sql := fmt.Sprintf("select %s from %s", leaveApplicationRows, m.table)
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

	sql += " Order by created_at desc "

	// 添加分页
	if data.Limit != -1 { // 约定 -1 表示查询全部
		sql += fmt.Sprintf(" LIMIT %d OFFSET %d", data.Limit, (data.Page-1)*data.Limit)
	}
	err = m.QueryRowsNoCacheCtx(ctx, &leaveApplication, sql, args...)

	switch {
	case err == nil:
		return leaveApplication, total, nil
	case errors.Is(err, sqlc.ErrNotFound):
		// 搜索时数据为空不是错误，返回空列表
		return []*LeaveApplication{}, total, nil
	default:
		return nil, 0, err
	}
}

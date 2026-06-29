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

var _ AttendanceReplenishModel = (*customAttendanceReplenishModel)(nil)

type (
	// AttendanceReplenishModel 补卡申请模型接口
	AttendanceReplenishModel interface {
		attendanceReplenishModel
		XUpdate(ctx context.Context, attendanceReplenish *AttendanceReplenish) error
		Search(ctx context.Context,
			data *types.SearchReplenishParams) ([]*AttendanceReplenish, int64, error)
	}

	customAttendanceReplenishModel struct {
		*defaultAttendanceReplenishModel
	}
)

// NewAttendanceReplenishModel 返回补卡申请模型
func NewAttendanceReplenishModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AttendanceReplenishModel {
	return &customAttendanceReplenishModel{
		defaultAttendanceReplenishModel: newAttendanceReplenishModel(conn, c, opts...),
	}
}

// XUpdate 动态更新补卡申请（只更新非零值字段）
func (m *customAttendanceReplenishModel) XUpdate(ctx context.Context, attendanceReplenish *AttendanceReplenish) error {
	var setClauses []string
	var args []interface{}

	// 补卡日期
	if !xtime.IsZeroTime(attendanceReplenish.OriginalDate) {
		setClauses = append(setClauses, "original_date = ?")
		args = append(args, attendanceReplenish.OriginalDate)
	}
	// 补卡类型
	if attendanceReplenish.ReplenishType != 0 {
		setClauses = append(setClauses, "replenish_type = ?")
		args = append(args, attendanceReplenish.ReplenishType)
	}
	// 补卡时间
	if !xtime.IsZeroTime(attendanceReplenish.ReplenishTime.Time) {
		setClauses = append(setClauses, "replenish_time = ?")
		args = append(args, attendanceReplenish.ReplenishTime)
	}
	// 补卡原因
	if attendanceReplenish.Reason != "" {
		setClauses = append(setClauses, "reason = ?")
		args = append(args, attendanceReplenish.Reason)
	}
	// 证明材料
	if attendanceReplenish.Evidence.String != "" {
		setClauses = append(setClauses, "evidence = ?")
		args = append(args, attendanceReplenish.Evidence.String)
	}
	// 状态
	if attendanceReplenish.Status != 0 {
		setClauses = append(setClauses, "status = ?")
		args = append(args, attendanceReplenish.Status)
	}
	// 审批人ID
	if attendanceReplenish.ApproverId.Int64 != 0 {
		setClauses = append(setClauses, "approver_id = ?")
		args = append(args, attendanceReplenish.ApproverId.Int64)
	}
	// 审批时间
	if !xtime.IsZeroTime(attendanceReplenish.ApproveTime.Time) {
		setClauses = append(setClauses, "approve_time = ?")
		args = append(args, attendanceReplenish.ApproveTime.Time)
	}
	// 审批意见
	if attendanceReplenish.ApproveRemark.String != "" {
		setClauses = append(setClauses, "approve_remark = ?")
		args = append(args, attendanceReplenish.ApproveRemark.String)
	}

	erpHrAttendanceReplenishIdKey := fmt.Sprintf("%s%v", cacheErpHrAttendanceReplenishIdPrefix, attendanceReplenish.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, strings.Join(setClauses, ", "))
		args = append(args, attendanceReplenish.Id)
		return conn.ExecCtx(ctx, query, args...)
	}, erpHrAttendanceReplenishIdKey)
	return err
}

// Search 搜索补卡申请
func (m *customAttendanceReplenishModel) Search(ctx context.Context,
	data *types.SearchReplenishParams) ([]*AttendanceReplenish, int64, error) {
	var attendanceReplenishes []*AttendanceReplenish

	conditions := []string{}
	args := []any{}

	// 按员工ID查询
	if data.EmployeeId != 0 {
		conditions = append(conditions, "employee_id = ?")
		args = append(args, data.EmployeeId)
	}
	// 补卡类型
	if data.ReplenishType != 0 {
		conditions = append(conditions, "replenish_type = ?")
		args = append(args, data.ReplenishType)
	}
	// 补卡原因模糊查询
	if data.Reason != "" {
		conditions = append(conditions, "reason LIKE ?")
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
	sql := fmt.Sprintf("select %s from %s", attendanceReplenishRows, m.table)
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

	sql += " Order by apply_time desc "

	// 添加分页
	if data.Limit != -1 { // 约定 -1 表示查询全部
		sql += fmt.Sprintf(" LIMIT %d OFFSET %d", data.Limit, (data.Page-1)*data.Limit)
	}
	err = m.QueryRowsNoCacheCtx(ctx, &attendanceReplenishes, sql, args...)

	switch {
	case err == nil:
		return attendanceReplenishes, total, nil
	case errors.Is(err, sqlc.ErrNotFound):
		// 搜索时数据为空不是错误，返回空列表
		return []*AttendanceReplenish{}, total, nil
	default:
		return nil, 0, err
	}
}

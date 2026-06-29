package model

import (
	"context"
	"database/sql"
	"erp/app/hr/rpc/internal/types"
	"erp/common/xtime"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AttendanceRecordModel = (*customAttendanceRecordModel)(nil)

type (
	// AttendanceRecordModel 考勤记录模型接口
	AttendanceRecordModel interface {
		attendanceRecordModel
		XUpdate(ctx context.Context, newData *AttendanceRecord) error
		Search(ctx context.Context, data *types.SearchAttendanceRecordParams) ([]*AttendanceRecord, int64, error)
		FindMissingClockOut(ctx context.Context, date time.Time) ([]*AttendanceRecord, error)
	}

	customAttendanceRecordModel struct {
		*defaultAttendanceRecordModel
	}
)

// NewAttendanceRecordModel 返回考勤记录模型
func NewAttendanceRecordModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AttendanceRecordModel {
	return &customAttendanceRecordModel{
		defaultAttendanceRecordModel: newAttendanceRecordModel(conn, c, opts...),
	}
}

// XUpdate 动态更新考勤记录（只更新非零值字段）
func (m *customAttendanceRecordModel) XUpdate(ctx context.Context, newData *AttendanceRecord) error {
	var setClauses []string
	var args []interface{}

	// 上班打卡时间
	if !newData.ClockIn.Time.IsZero() {
		setClauses = append(setClauses, "clock_in = ?")
		args = append(args, newData.ClockIn.Time)
	}
	// 下班打卡时间
	if !xtime.IsZeroTime(newData.ClockOut.Time) {
		setClauses = append(setClauses, "clock_out = ?")
		args = append(args, newData.ClockOut.Time)
	}

	// 布尔状态字段（新版表结构）
	// 注意：布尔值需要显式设置，这里通过传入值判断是否需要更新
	if newData.IsAmMissing != 0 {
		setClauses = append(setClauses, "is_am_missing = ?")
		args = append(args, newData.IsAmMissing == 1)
	}
	if newData.IsLate != 0 {
		setClauses = append(setClauses, "is_late = ?")
		args = append(args, newData.IsLate == 1)
	}
	if newData.IsPmMissing != 0 {
		setClauses = append(setClauses, "is_pm_missing = ?")
		args = append(args, newData.IsPmMissing == 1)
	}
	if newData.IsEarlyLeave != 0 {
		setClauses = append(setClauses, "is_early_leave = ?")
		args = append(args, newData.IsEarlyLeave == 1)
	}

	// 工作时长
	if newData.WorkHours != 0 {
		setClauses = append(setClauses, "work_hours = ?")
		args = append(args, newData.WorkHours)
	}
	// 加班时长
	if newData.OvertimeHours != 0 {
		setClauses = append(setClauses, "overtime_hours = ?")
		args = append(args, newData.OvertimeHours)
	}
	// 备注
	if newData.Remark.String != "" {
		setClauses = append(setClauses, "remark = ?")
		args = append(args, newData.Remark.String)
	}

	erpHrAttendanceRecordIdKey := fmt.Sprintf("%s%v", cacheErpHrAttendanceRecordIdPrefix, newData.Id)
	erpHrAttendanceRecordEmployeeIdDateKey := fmt.Sprintf("%s%v:%v", cacheErpHrAttendanceRecordEmployeeIdDatePrefix, newData.EmployeeId, newData.Date)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, strings.Join(setClauses, ", "))
		args = append(args, newData.Id)
		return conn.ExecCtx(ctx, query, args...)
	}, erpHrAttendanceRecordIdKey, erpHrAttendanceRecordEmployeeIdDateKey)
	return err
}

// Search 搜索考勤记录
func (m *customAttendanceRecordModel) Search(ctx context.Context,
	data *types.SearchAttendanceRecordParams) ([]*AttendanceRecord, int64, error) {
	var attendanceRecords []*AttendanceRecord

	conditions := []string{}
	args := []any{}

	// 按员工ID查询
	if data.EmployeeId != 0 {
		conditions = append(conditions, "employee_id = ?")
		args = append(args, data.EmployeeId)
	}
	// 日期范围
	if !xtime.IsZeroTime(data.StartDate) {
		conditions = append(conditions, "date >= ?")
		args = append(args, data.StartDate)
	}
	if !xtime.IsZeroTime(data.EndDate) {
		conditions = append(conditions, "date <= ?")
		args = append(args, data.EndDate)
	}
	// 迟到筛选
	if data.IsLate {
		conditions = append(conditions, "is_late = ?")
		args = append(args, true)
	}
	// 早退筛选
	if data.IsEarlyLeave {
		conditions = append(conditions, "is_early_leave = ?")
		args = append(args, true)
	}
	// 缺上午卡筛选
	if data.IsAmMissing {
		conditions = append(conditions, "is_am_missing = ?")
		args = append(args, true)
	}
	// 缺下午卡筛选
	if data.IsPmMissing {
		conditions = append(conditions, "is_pm_missing = ?")
		args = append(args, true)
	}
	// 备注模糊查询
	if data.Remark != "" {
		conditions = append(conditions, "remark LIKE ?")
		args = append(args, "%"+data.Remark+"%")
	}

	// 构建完整 SQL
	sql := fmt.Sprintf("select %s from %s", attendanceRecordRows, m.table)
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

	sql += " Order by date desc "

	// 添加分页
	if data.Limit != -1 { // 约定 -1 表示查询全部
		sql += fmt.Sprintf(" LIMIT %d OFFSET %d", data.Limit, (data.Page-1)*data.Limit)
	}
	err = m.QueryRowsNoCacheCtx(ctx, &attendanceRecords, sql, args...)

	switch {
	case err == nil:
		return attendanceRecords, total, nil
	case errors.Is(err, sqlc.ErrNotFound):
		// 搜索时数据为空不是错误，返回空列表
		return []*AttendanceRecord{}, total, nil
	default:
		return nil, 0, err
	}
}

// FindMissingClockOut 查找缺少下班打卡的记录
func (m *customAttendanceRecordModel) FindMissingClockOut(ctx context.Context, date time.Time) ([]*AttendanceRecord, error) {
	var attendanceRecords []*AttendanceRecord
	query := fmt.Sprintf("select %s from %s where `date` = ? and `clock_out` is null", attendanceRecordRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &attendanceRecords, query, date)
	switch {
	case err == nil:
		return attendanceRecords, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

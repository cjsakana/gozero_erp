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

var _ PayrollRecordModel = (*customPayrollRecordModel)(nil)

type (
	// PayrollRecordModel 薪资记录模型接口
	PayrollRecordModel interface {
		payrollRecordModel
		BulkInsert(data []*PayrollRecord) ([]*types.BulkAddPayrollRecordErrItem, error)
		XUpdate(ctx context.Context, newData *PayrollRecord) error
		Search(ctx context.Context, data *types.SearchPayrollRecordParams) ([]*PayrollRecord, int64, error)
	}

	customPayrollRecordModel struct {
		*defaultPayrollRecordModel
		conn sqlx.SqlConn
	}
)

// NewPayrollRecordModel 返回薪资记录模型
func NewPayrollRecordModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PayrollRecordModel {
	return &customPayrollRecordModel{
		defaultPayrollRecordModel: newPayrollRecordModel(conn, c, opts...),
		conn:                      conn,
	}
}

// XUpdate 动态更新薪资记录（只更新非零值字段）
func (m *customPayrollRecordModel) XUpdate(ctx context.Context, newData *PayrollRecord) error {
	var setClauses []string
	var args []interface{}

	// 奖金
	if newData.Bonus != 0 {
		setClauses = append(setClauses, "bonus = ?")
		args = append(args, newData.Bonus)
	}
	// 扣款
	if newData.Deductions != 0 {
		setClauses = append(setClauses, "deductions = ?")
		args = append(args, newData.Deductions)
	}
	// 实发工资
	if newData.NetSalary.Float64 != 0 {
		setClauses = append(setClauses, "net_salary = ?")
		args = append(args, newData.NetSalary.Float64)
	}
	// 核算人ID
	if newData.CalculatedBy.Int64 != 0 {
		setClauses = append(setClauses, "calculated_by = ?")
		args = append(args, newData.CalculatedBy.Int64)
	}
	// 核算时间
	if !xtime.IsZeroTime(newData.CalculatedAt.Time) {
		setClauses = append(setClauses, "calculated_at = ?")
		args = append(args, newData.CalculatedAt.Time)
	}
	// 状态
	if newData.Status != 0 {
		setClauses = append(setClauses, "status = ?")
		args = append(args, newData.Status)
	}
	// 描述
	if newData.Description.String != "" {
		setClauses = append(setClauses, "description = ?")
		args = append(args, newData.Description.String)
	}
	// 发放时间
	if !xtime.IsZeroTime(newData.PaymentAt.Time) {
		setClauses = append(setClauses, "payment_at = ?")
		args = append(args, newData.PaymentAt.Time)
	}

	erpHrPayrollRecordIdKey := fmt.Sprintf("%s%v", cacheErpHrPayrollRecordIdPrefix, newData.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, strings.Join(setClauses, ", "))
		args = append(args, newData.Id)
		return conn.ExecCtx(ctx, query, args...)
	}, erpHrPayrollRecordIdKey)
	return err
}

// Search 搜索薪资记录
func (m *customPayrollRecordModel) Search(ctx context.Context, data *types.SearchPayrollRecordParams) ([]*PayrollRecord, int64, error) {
	var records []*PayrollRecord

	conditions := []string{}
	args := []any{}

	// 按员工ID查询
	if data.EmployeeId != 0 {
		conditions = append(conditions, "employee_id = ?")
		args = append(args, data.EmployeeId)
	}
	// 状态
	if data.Status != 0 {
		conditions = append(conditions, "status = ?")
		args = append(args, data.Status)
	}
	// 描述模糊查询
	if data.Description != "" {
		conditions = append(conditions, "description LIKE ?")
		args = append(args, "%"+data.Description+"%")
	}
	// 核算人ID
	if data.CalculatedBy != 0 {
		conditions = append(conditions, "calculated_by = ?")
		args = append(args, data.CalculatedBy)
	}
	// 核算日期范围
	if !xtime.IsZeroTime(data.StartCalculatedDate) {
		conditions = append(conditions, "calculated_at >= ?")
		args = append(args, data.StartCalculatedDate)
	}
	if !xtime.IsZeroTime(data.EndCalculatedDate) {
		conditions = append(conditions, "calculated_at <= ?")
		args = append(args, data.EndCalculatedDate)
	}
	// 发放日期范围
	if !xtime.IsZeroTime(data.StartPaymentDate) {
		conditions = append(conditions, "payment_at >= ?")
		args = append(args, data.StartPaymentDate)
	}
	if !xtime.IsZeroTime(data.EndPaymentDate) {
		conditions = append(conditions, "payment_at <= ?")
		args = append(args, data.EndPaymentDate)
	}
	// 薪资月份
	if !xtime.IsZeroTime(data.PaymentMonth) {
		conditions = append(conditions, "payment_month = ?")
		args = append(args, data.PaymentMonth)
	}

	// 构建完整 SQL
	sql := fmt.Sprintf("select %s from %s", payrollRecordRows, m.table)
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
	err = m.QueryRowsNoCacheCtx(ctx, &records, sql, args...)

	switch {
	case err == nil:
		return records, total, nil
	case errors.Is(err, sqlc.ErrNotFound):
		// 搜索时数据为空不是错误，返回空列表
		return []*PayrollRecord{}, total, nil
	default:
		return nil, 0, err
	}
}

// BulkInsert 批量插入薪资记录
func (m *customPayrollRecordModel) BulkInsert(data []*PayrollRecord) ([]*types.BulkAddPayrollRecordErrItem, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, payrollRecordRowsExpectAutoSet)

	blk, err := sqlx.NewBulkInserter(m.conn, query)
	if err != nil {
		panic(err)
	}
	defer blk.Flush()

	results := make([]*types.BulkAddPayrollRecordErrItem, 0, len(data))

	for _, v := range data {
		item := &types.BulkAddPayrollRecordErrItem{
			EmployeeId: v.EmployeeId,
			Success:    true,
		}
		err := blk.Insert(v.Id, v.EmployeeId, v.PaymentMonth, v.BaseSalary, v.Bonus, v.Deductions,
			v.NetSalary, v.CalculatedBy, v.CalculatedAt, v.Status, v.Description, v.PaymentAt)
		if err != nil {
			item.Success = false
			item.Err = err
		}

		results = append(results, item)
	}
	return results, nil
}

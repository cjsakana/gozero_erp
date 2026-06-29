package model

import (
	"context"
	"database/sql"
	"github.com/zeromicro/go-zero/core/stringx"
	"time"
	types2 "erp/app/hr/rpc/internal/types"
	"erp/common/xtime"
	"errors"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EmployeeDetailModel = (*customEmployeeDetailModel)(nil)

type (
	// EmployeeDetailModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEmployeeDetailModel.
	EmployeeDetailModel interface {
		employeeDetailModel
		XUpdate(ctx context.Context, employeeDetail *EmployeeDetail) error
		Search(ctx context.Context, data *types2.SearchEmployeeDetailParam) ([]*EmployeeDetail, int64, error)
		BulkInsert(data []*EmployeeDetail) ([]*types2.BulkInsertResult, error)
		ClearLeaveDate(ctx context.Context, employeeNo string) error
		SetLeaveDate(ctx context.Context, employeeNo string, t time.Time) error
	}

	customEmployeeDetailModel struct {
		*defaultEmployeeDetailModel
		conn sqlx.SqlConn
	}
)

// NewEmployeeDetailModel returns a model for the database table.
func NewEmployeeDetailModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) EmployeeDetailModel {
	return &customEmployeeDetailModel{
		defaultEmployeeDetailModel: newEmployeeDetailModel(conn, c, opts...),
		conn:                       conn,
	}
}
func (m *customEmployeeDetailModel) XUpdate(ctx context.Context, employeeDetail *EmployeeDetail) error {
	var setClauses []string
	var args []interface{}

	if employeeDetail.Account.String != "" {
		setClauses = append(setClauses, "account = ?")
		args = append(args, employeeDetail.Account.String)
	}
	if employeeDetail.DepartmentId != 0 {
		setClauses = append(setClauses, "department_id = ?")
		args = append(args, employeeDetail.DepartmentId)
	}
	if employeeDetail.PositionId != 0 {
		setClauses = append(setClauses, "position_id = ?")
		args = append(args, employeeDetail.PositionId)
	}
	if employeeDetail.Salary.Float64 != 0 {
		setClauses = append(setClauses, "salary = ?")
		args = append(args, employeeDetail.Salary.Float64)
	}
	if !xtime.IsZeroTime(employeeDetail.HireDate) {
		setClauses = append(setClauses, "hire_date = ?")
		args = append(args, employeeDetail.HireDate)
	}
	if !xtime.IsZeroTime(employeeDetail.LeaveDate.Time) {
		setClauses = append(setClauses, "leave_date = ?")
		args = append(args, employeeDetail.LeaveDate.Time)
	}
	if employeeDetail.Name != "" {
		setClauses = append(setClauses, "name = ?")
		args = append(args, employeeDetail.Name)
	}

	erpHrEmployeeDetailEmployeeNoKey := fmt.Sprintf("%s%v", cacheErpHrEmployeeDetailEmployeeNoPrefix, employeeDetail.EmployeeNo)
	erpHrEmployeeDetailIdKey := fmt.Sprintf("%s%v", cacheErpHrEmployeeDetailIdPrefix, employeeDetail.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ? and `leave_date` is null ", m.table, strings.Join(setClauses, ", "))
		args = append(args, employeeDetail.Id)
		return conn.ExecCtx(ctx, query, args...)
	}, erpHrEmployeeDetailEmployeeNoKey, erpHrEmployeeDetailIdKey)
	return err
}

func (m *customEmployeeDetailModel) Search(ctx context.Context, data *types2.SearchEmployeeDetailParam) ([]*EmployeeDetail, int64, error) {
	var employeeDetails []*EmployeeDetail

	conditions := []string{}
	args := []any{}

	if data.Gender != 0 {
		conditions = append(conditions, "gender = ?")
		args = append(args, data.Gender)
	}

	if data.DepartmentId != 0 {
		conditions = append(conditions, "department_id = ?")
		args = append(args, data.DepartmentId)
	}
	if data.PositionId != 0 {
		conditions = append(conditions, "position_id = ?")
		args = append(args, data.PositionId)
	}
	if data.Salary != 0 {
		conditions = append(conditions, "salary = ?")
		args = append(args, data.Salary)
	}
	if !xtime.IsZeroTime(data.HireDate) {
		conditions = append(conditions, "hire_date = ?")
		args = append(args, data.HireDate)
	}
	if data.Name != "" {
		conditions = append(conditions, "name like ?")
		args = append(args, "%"+data.Name+"%")
	}
	if data.Resigned == 1 {
		conditions = append(conditions, "leave_date is not null")
	} else {
		conditions = append(conditions, "leave_date is null")
	}

	// 构建完整 SQL
	sql := fmt.Sprintf("select %s from %s", employeeDetailRows, m.table)
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

	sql += " Order by hire_date desc "

	// 添加分页
	if data.Limit != -1 { // 约定 -1 表示查询全部
		sql += fmt.Sprintf(" LIMIT %d OFFSET %d", data.Limit, (data.Page-1)*data.Limit)
	}
	err = m.QueryRowsNoCacheCtx(ctx, &employeeDetails, sql, args...)

	switch {
	case err == nil:
		return employeeDetails, total, nil
	case errors.Is(err, sqlc.ErrNotFound):
		// 搜索时数据为空不是错误，返回空列表
		return []*EmployeeDetail{}, total, nil
	default:
		return nil, 0, err
	}
}

func (m *customEmployeeDetailModel) BulkInsert(data []*EmployeeDetail) ([]*types2.BulkInsertResult, error) {
	xEmployeeDetailRowsExpectAutoSet := strings.Join(stringx.Remove(employeeDetailFieldNames, "`leave_date`"), ",")

	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?,  ?, ?, ?, ?, ?, ?)", m.table, xEmployeeDetailRowsExpectAutoSet)

	blk, err := sqlx.NewBulkInserter(m.conn, query)
	if err != nil {
		panic(err)
	}
	defer blk.Flush()

	results := make([]*types2.BulkInsertResult, len(data))

	// 临时索引映射：每条 Insert 对应 data 的索引
	batchIndices := []int{}

	blk.SetResultHandler(func(result sql.Result, e error) {
		// 复制当前批次索引，避免闭包持有同一个底层数组被后续修改
		localBatch := append([]int(nil), batchIndices...)
		if e != nil {
			// 批量执行失败，则这一批全失败
			for _, idx := range localBatch {
				results[idx] = &types2.BulkInsertResult{Index: idx, Success: false, Err: errors.New(e.Error())}
			}
		} else {
			// 成功时 RowsAffected 对应批量插入行数
			rowsAffected, _ := result.RowsAffected()
			for i, idx := range localBatch {
				success := i < int(rowsAffected) // 假设前 rowsAffected 条成功
				results[idx] = &types2.BulkInsertResult{Index: idx, Success: success, Err: nil}
			}
		}
		batchIndices = []int{} // 清空批次
	})

	for i, v := range data {
		// 添加到当前批次索引
		batchIndices = append(batchIndices, i)

		if err := blk.Insert(
			v.Id, v.EmployeeNo, v.Name, v.IdCard, v.Account.String, v.Gender,
			v.BirthDate.Format("2006-01-02 15:04:05"), v.DepartmentId,
			v.PositionId, v.Salary.Float64, v.HireDate.Format("2006-01-02 15:04:05"),
		); err != nil {
			// 如果 Insert 参数错误（立即报错），记录失败
			results[i] = &types2.BulkInsertResult{Index: i, Success: false, Err: errors.New(err.Error())}
			batchIndices = batchIndices[:len(batchIndices)-1] // 从批次移除
		}
	}

	// 最后 Flush 会触发 SetResultHandler
	return results, nil
}

// ClearLeaveDate 将指定员工的离职日期清空为 NULL
func (m *customEmployeeDetailModel) ClearLeaveDate(ctx context.Context, employeeNo string) error {
	erpHrEmployeeDetailEmployeeNoKey := fmt.Sprintf("%s%v", cacheErpHrEmployeeDetailEmployeeNoPrefix, employeeNo)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set leave_date = null where `employee_no` = ?", m.table)
		return conn.ExecCtx(ctx, query, employeeNo)
	}, erpHrEmployeeDetailEmployeeNoKey)
	return err
}

// SetLeaveDate 设置指定员工的离职日期
func (m *customEmployeeDetailModel) SetLeaveDate(ctx context.Context, employeeNo string, t time.Time) error {
	erpHrEmployeeDetailEmployeeNoKey := fmt.Sprintf("%s%v", cacheErpHrEmployeeDetailEmployeeNoPrefix, employeeNo)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set leave_date = ? where `employee_no` = ?", m.table)
		return conn.ExecCtx(ctx, query, t, employeeNo)
	}, erpHrEmployeeDetailEmployeeNoKey)
	return err
}

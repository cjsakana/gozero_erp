package model

import (
	"context"
	"database/sql"
	types2 "erp/app/hr/rpc/internal/types"
	"errors"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ DepartmentModel = (*customDepartmentModel)(nil)

type (
	// DepartmentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDepartmentModel.
	DepartmentModel interface {
		departmentModel
		XUpdate(ctx context.Context, department *Department) error
		Search(ctx context.Context, data *types2.SearchDepartmentParams) ([]*Department, int64, error)
		//BulkInsert(data []*Department) ([]*types2.BulkInsertResult, error)
	}

	customDepartmentModel struct {
		*defaultDepartmentModel
		conn sqlx.SqlConn
	}
)

// NewDepartmentModel returns a model for the database table.
func NewDepartmentModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) DepartmentModel {
	return &customDepartmentModel{
		defaultDepartmentModel: newDepartmentModel(conn, c, opts...),
		conn:                   conn,
	}
}

func (m *customDepartmentModel) XUpdate(ctx context.Context, department *Department) error {
	var setClauses []string
	var args []interface{}

	if department.Name != "" {
		setClauses = append(setClauses, "name = ?")
		args = append(args, department.Name)
	}
	if department.ParentId.Int64 != 0 {
		setClauses = append(setClauses, "parent_id = ?")
		args = append(args, department.ParentId.Int64)
	}
	if department.Code.String != "" {
		setClauses = append(setClauses, "code = ?")
		args = append(args, department.Code.String)
	}
	if department.ManagerNo.String != "" {
		setClauses = append(setClauses, "manager_no = ?")
		args = append(args, department.ManagerNo.String)
	}
	if department.ManagerName != "" {
		setClauses = append(setClauses, "manager_name = ?")
		args = append(args, department.ManagerName)
	}

	erpHrDepartmentIdKey := fmt.Sprintf("%s%v", cacheErpHrDepartmentIdPrefix, department.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, strings.Join(setClauses, ", "))
		args = append(args, department.Id)
		return conn.ExecCtx(ctx, query, args...)
	}, erpHrDepartmentIdKey)
	return err
}

func (m *customDepartmentModel) Search(ctx context.Context, data *types2.SearchDepartmentParams) ([]*Department, int64, error) {
	var departments []*Department

	conditions := []string{}
	args := []any{}

	if data.Name != "" {
		conditions = append(conditions, "name LIKE ?")
		args = append(args, "%"+data.Name+"%")
	}
	if data.ParentId != 0 {
		conditions = append(conditions, "parent_id = ?")
		args = append(args, data.ParentId)
	}
	if data.Code != "" {
		conditions = append(conditions, "code like ?")
		args = append(args, "%"+data.Code+"%")
	}

	// 构建完整 SQL
	sql := fmt.Sprintf("select %s from %s", departmentRows, m.table)
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

	// 添加分页
	if data.Limit != -1 { // 约定 -1 表示查询全部
		sql += fmt.Sprintf(" LIMIT %d OFFSET %d", data.Limit, (data.Page-1)*data.Limit)
	}
	err = m.QueryRowsNoCacheCtx(ctx, &departments, sql, args...)

	switch {
	case err == nil:
		return departments, total, nil
	case errors.Is(err, sqlc.ErrNotFound):
		// 搜索时数据为空不是错误，返回空列表
		return []*Department{}, total, nil
	default:
		return nil, 0, err
	}
}

//func (m *customDepartmentModel) BulkInsert(data []*Department) ([]*types2.BulkInsertResult, error) {
//	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, departmentRowsExpectAutoSet)
//
//	blk, err := sqlx.NewBulkInserter(m.conn, query)
//	if err != nil {
//		panic(err)
//	}
//	defer blk.Flush()
//	results := make([]*types2.BulkInsertResult, len(data))
//
//	// 临时索引映射：每条 Insert 对应 data 的索引
//	batchIndices := []int{}
//	blk.SetResultHandler(func(result sql.Result, e error) {
//		if e != nil {
//			// 批量执行失败，则这一批全失败
//			for _, idx := range batchIndices {
//				results[idx] = &types2.BulkInsertResult{Index: idx, Success: false, Err: e}
//			}
//		} else {
//			// 成功时 RowsAffected 对应批量插入行数
//			rowsAffected, _ := result.RowsAffected()
//			for i, idx := range batchIndices {
//				success := i < int(rowsAffected) // 假设前 rowsAffected 条成功
//				results[idx] = &types2.BulkInsertResult{Index: idx, Success: success, Err: nil}
//			}
//		}
//		batchIndices = []int{} // 清空批次
//	})
//
//	for i, v := range data {
//		// 添加到当前批次索引
//		batchIndices = append(batchIndices, i)
//		err := blk.Insert(v.Name, v.ParentId, v.Code, v.ManagerNo, v.ManagerName)
//		if err != nil {
//			// 如果 Insert 参数错误（立即报错），记录失败
//			results[i] = &types2.BulkInsertResult{Index: i, Success: false, Err: err}
//			batchIndices = batchIndices[:len(batchIndices)-1] // 从批次移除
//		}
//	}
//	return results, nil
//}

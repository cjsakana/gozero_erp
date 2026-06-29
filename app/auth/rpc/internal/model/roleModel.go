package model

import (
	"context"
	"database/sql"
	"erp/app/auth/rpc/internal/types"
	"errors"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RoleModel = (*customRoleModel)(nil)

type (
	// RoleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRoleModel.
	RoleModel interface {
		roleModel
		XUpdate(ctx context.Context, newData *Role) error
		SearchRoles(ctx context.Context, data *types.SearchRole) ([]*Role, int64, error)
		XDeleteTX(ctx context.Context, id int64) ([]int64, error)
	}

	customRoleModel struct {
		*defaultRoleModel
		conn sqlx.SqlConn
	}
)

// NewRoleModel returns a model for the database table.
func NewRoleModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) RoleModel {
	return &customRoleModel{
		defaultRoleModel: newRoleModel(conn, c, opts...),
		conn:             conn,
	}
}

func (m *customRoleModel) XUpdate(ctx context.Context, newData *Role) error {
	var setClauses []string
	var args []interface{}

	if newData.Code != "" {
		setClauses = append(setClauses, "code = ?")
		args = append(args, newData.Code)
	}
	if newData.Name != "" {
		setClauses = append(setClauses, "name = ?")
		args = append(args, newData.Name)
	}
	if newData.Description.String != "" {
		setClauses = append(setClauses, "description = ?")
		args = append(args, newData.Description.String)
	}

	erpEmployeesEmployeeIdKey := fmt.Sprintf("%s%v", cacheErpAuthRoleIdPrefix, newData.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, strings.Join(setClauses, ", "))
		args = append(args, newData.Id)
		return conn.ExecCtx(ctx, query, args...)
	}, erpEmployeesEmployeeIdKey)
	return err
}

func (m *customRoleModel) SearchRoles(ctx context.Context, data *types.SearchRole) ([]*Role, int64, error) {
	var roles []*Role

	conditions := []string{}
	args := []any{}

	if data.Code != "" {
		conditions = append(conditions, "code LIKE ?")
		args = append(args, "%"+data.Code+"%")
	}

	if data.Name != "" {
		conditions = append(conditions, "name LIKE ?")
		args = append(args, "%"+data.Name+"%")
	}

	if data.Description != "" {
		conditions = append(conditions, "description like ?")
		args = append(args, data.Description)
	}

	// 构建完整 SQL
	sql := fmt.Sprintf("select %s from %s", roleRows, m.table)
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
	err = m.QueryRowsNoCacheCtx(ctx, &roles, sql, args...)

	switch {
	case err == nil:
		return roles, total, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, 0, ErrNotFound
	default:
		return nil, 0, err
	}
}

func (m *customRoleModel) XDeleteTX(ctx context.Context, id int64) ([]int64, error) {

	rpIds := []int64{}
	err := m.QueryRowsNoCacheCtx(ctx, &rpIds, "select id from role_permission where role_id = ?", id)
	if err != nil {
		return nil, err
	}
	err = m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		_, err := session.ExecCtx(ctx, query, id)
		if err != nil {
			return err
		}
		for _, rpId := range rpIds {
			_, err := session.ExecCtx(ctx, "delete from role_permission where `id` = ?", rpId)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return rpIds, err
}

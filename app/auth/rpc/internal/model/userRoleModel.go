package model

import (
	"context"
	"database/sql"
	"erp/app/auth/rpc/internal/types"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
)

var _ UserRoleModel = (*customUserRoleModel)(nil)

type (
	// UserRoleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserRoleModel.
	UserRoleModel interface {
		userRoleModel
		FindRolesByUserId(ctx context.Context, userId int64) ([]*UserRole, error)
		XUpdate(ctx context.Context, id, roleId, userId int64) error
		SearchUserRoles(ctx context.Context, data *types.SearchUserRoleParams) ([]*UserRole, int64, error)
		DeleteByUserId(ctx context.Context, userId int64) ([]int64, error)
	}

	customUserRoleModel struct {
		*defaultUserRoleModel
		conn sqlx.SqlConn
	}
)

// NewUserRoleModel returns a model for the database table.
func NewUserRoleModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserRoleModel {
	return &customUserRoleModel{
		defaultUserRoleModel: newUserRoleModel(conn, c, opts...),
		conn:                 conn,
	}
}

func (m *customUserRoleModel) FindRolesByUserId(ctx context.Context, userId int64) ([]*UserRole, error) {
	var resp []*UserRole
	query := fmt.Sprintf("select %s from %s where `user_id` = ?", userRoleRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)

	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUserRoleModel) XUpdate(ctx context.Context, id, roleId, userId int64) error {
	roles, err := m.FindRolesByUserId(ctx, roleId)
	if err != nil {
		return err
	}
	for _, role := range roles {
		if role.Id == roleId {
			return err
		}
	}
	err = m.Update(ctx, &UserRole{
		Id:     id,
		UserId: sql.NullInt64{Int64: userId, Valid: true},
		RoleId: sql.NullInt64{Int64: roleId, Valid: true},
	})
	if err != nil {
		return err
	}
	return nil
}

func (m *customUserRoleModel) SearchUserRoles(ctx context.Context, data *types.SearchUserRoleParams) ([]*UserRole, int64, error) {
	var userRoles []*UserRole

	conditions := []string{}
	args := []any{}

	if data.UserId != 0 {
		conditions = append(conditions, "user_id = ?")
		args = append(args, data.UserId)
	}
	if data.RoleId != 0 {
		conditions = append(conditions, "role_id = ?")
		args = append(args, data.RoleId)
	}

	// 构建完整 SQL
	sql := fmt.Sprintf("select %s from %s", userRoleRows, m.table)
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
	err = m.QueryRowsNoCacheCtx(ctx, &userRoles, sql, args...)

	switch {
	case err == nil:
		return userRoles, total, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, 0, ErrNotFound
	default:
		return nil, 0, err
	}
}
func (m *customUserRoleModel) DeleteByUserId(ctx context.Context, userId int64) ([]int64, error) {
	var urs []*UserRole
	query := fmt.Sprintf("delete from %s where `user_id`=?", m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &urs, query, userId)
	if err != nil {
		return nil, err
	}
	var ids []int64
	err = m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {

		for _, ur := range urs {
			_, err := session.ExecCtx(ctx, fmt.Sprintf("delete from %s where `id`=?", m.table), ur.Id)
			if err != nil {
				return err
			}
			ids = append(ids, ur.Id)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return ids, nil
}

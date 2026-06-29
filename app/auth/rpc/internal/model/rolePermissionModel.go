package model

import (
	"context"
	"database/sql"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RolePermissionModel = (*customRolePermissionModel)(nil)

type (
	// PermissionWithRolePermissionId 包含role_permission_id的权限信息
	PermissionWithRolePermissionId struct {
		RolePermissionId int64          `db:"role_permission_id"` // role_permission表的主键ID
		Id               int64          `db:"id"`                 // 权限ID
		ParentId         int64          `db:"parent_id"`          // 父级权限ID。0：根节点
		Code             sql.NullString `db:"code"`               // 权限代码
		Description      sql.NullString `db:"description"`        // 权限描述/名称
		Url              sql.NullString `db:"url"`
		Method           sql.NullString `db:"method"`
	}

	// RolePermissionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRolePermissionModel.
	RolePermissionModel interface {
		rolePermissionModel
		FindPermissionsByRoleId(ctx context.Context, roleId int64) ([]*Permission, error)
		FindPermissionsWithRolePermissionIdByRoleId(ctx context.Context, roleId int64) ([]*PermissionWithRolePermissionId, error)
	}

	customRolePermissionModel struct {
		*defaultRolePermissionModel
	}
)

// NewRolePermissionModel returns a model for the database table.
func NewRolePermissionModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) RolePermissionModel {
	return &customRolePermissionModel{
		defaultRolePermissionModel: newRolePermissionModel(conn, c, opts...),
	}
}

func (m *customRolePermissionModel) FindPermissionsByRoleId(ctx context.Context, roleId int64) ([]*Permission, error) {
	var resp []*Permission
	query := `SELECT p.id, p.parent_id, p.code, p.description, p.url, p.method
		FROM permission p
		INNER JOIN role_permission rp ON p.id = rp.permission_id
		WHERE rp.role_id = ?;
		`
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, roleId)

	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customRolePermissionModel) FindPermissionsWithRolePermissionIdByRoleId(ctx context.Context, roleId int64) ([]*PermissionWithRolePermissionId, error) {
	var resp []*PermissionWithRolePermissionId
	query := `SELECT rp.id as role_permission_id, p.id, p.parent_id, p.code, p.description, p.url, p.method
		FROM permission p
		INNER JOIN role_permission rp ON p.id = rp.permission_id
		WHERE rp.role_id = ?;
		`
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, roleId)

	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

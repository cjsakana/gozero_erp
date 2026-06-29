package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PermissionModel = (*customPermissionModel)(nil)

type (
	// PermissionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPermissionModel.
	PermissionModel interface {
		permissionModel
		FindAll(ctx context.Context) ([]*Permission, error)
	}

	customPermissionModel struct {
		*defaultPermissionModel
	}
)

// NewPermissionModel returns a model for the database table.
func NewPermissionModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PermissionModel {
	return &customPermissionModel{
		defaultPermissionModel: newPermissionModel(conn, c, opts...),
	}
}

func (m *customPermissionModel) FindAll(ctx context.Context) ([]*Permission, error) {
	query := fmt.Sprintf("select %s from %s", permissionRows, m.table)
	var permissions []*Permission
	err := m.QueryRowsNoCacheCtx(ctx, &permissions, query)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

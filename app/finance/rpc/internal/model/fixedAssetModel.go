package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FixedAssetModel = (*customFixedAssetModel)(nil)

type (
	// FixedAssetModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFixedAssetModel.
	FixedAssetModel interface {
		fixedAssetModel
		Search(ctx context.Context, assetNo, assetName, category string, supplierId, departmentId, status int64, page, limit int64) ([]*FixedAsset, int64, error)
	}

	customFixedAssetModel struct {
		*defaultFixedAssetModel
		conn sqlx.SqlConn
	}
)

// NewFixedAssetModel returns a model for the database table.
func NewFixedAssetModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FixedAssetModel {
	return &customFixedAssetModel{
		defaultFixedAssetModel: newFixedAssetModel(conn, c, opts...),
		conn:                   conn,
	}
}

func (m *customFixedAssetModel) Search(ctx context.Context, assetNo, assetName, category string, supplierId, departmentId, status int64, page, limit int64) ([]*FixedAsset, int64, error) {
	var fixedAssets []*FixedAsset
	conditions := []string{}
	args := []any{}

	if assetNo != "" {
		conditions = append(conditions, "asset_no = ?")
		args = append(args, assetNo)
	}
	if assetName != "" {
		conditions = append(conditions, "asset_name LIKE ?")
		args = append(args, "%"+assetName+"%")
	}
	if category != "" {
		conditions = append(conditions, "category = ?")
		args = append(args, category)
	}
	if supplierId != 0 {
		conditions = append(conditions, "supplier_id = ?")
		args = append(args, supplierId)
	}
	if departmentId != 0 {
		conditions = append(conditions, "department_id = ?")
		args = append(args, departmentId)
	}
	if status != 0 {
		conditions = append(conditions, "status = ?")
		args = append(args, status)
	}

	base := fmt.Sprintf("select %s from %s", fixedAssetRows, m.table)
	countBase := fmt.Sprintf("select count(*) from %s", m.table)
	if len(conditions) > 0 {
		where := " where " + strings.Join(conditions, " AND ")
		base += where
		countBase += where
	}

	var total int64
	if err := m.QueryRowNoCacheCtx(ctx, &total, countBase, args...); err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		base += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, (page-1)*limit)
	}

	if err := m.QueryRowsNoCacheCtx(ctx, &fixedAssets, base, args...); err != nil {
		return nil, 0, err
	}
	return fixedAssets, total, nil
}

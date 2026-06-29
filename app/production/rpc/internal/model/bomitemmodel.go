package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BomItemModel = (*customBomItemModel)(nil)

type (
	// BomItemModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBomItemModel.
	BomItemModel interface {
		bomItemModel
		FindByBomId(ctx context.Context, bomId int64) ([]*BomItem, error)
	}

	customBomItemModel struct {
		*defaultBomItemModel
	}
)

// NewBomItemModel returns a model for the database table.
func NewBomItemModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) BomItemModel {
	return &customBomItemModel{
		defaultBomItemModel: newBomItemModel(conn, c, opts...),
	}
}

func (m *customBomItemModel) FindByBomId(ctx context.Context, bomId int64) ([]*BomItem, error) {
	var items []*BomItem
	query := "SELECT * FROM bom_item WHERE bom_id = ?"
	err := m.QueryRowsNoCacheCtx(ctx, &items, query, bomId)
	return items, err
}

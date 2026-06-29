package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MaterialRequisitionItemModel = (*customMaterialRequisitionItemModel)(nil)

type (
	// MaterialRequisitionItemModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMaterialRequisitionItemModel.
	MaterialRequisitionItemModel interface {
		materialRequisitionItemModel
		FindByRequisitionId(ctx context.Context, requisitionId int64) ([]*MaterialRequisitionItem, error)
	}

	customMaterialRequisitionItemModel struct {
		*defaultMaterialRequisitionItemModel
	}
)

// NewMaterialRequisitionItemModel returns a model for the database table.
func NewMaterialRequisitionItemModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) MaterialRequisitionItemModel {
	return &customMaterialRequisitionItemModel{
		defaultMaterialRequisitionItemModel: newMaterialRequisitionItemModel(conn, c, opts...),
	}
}

func (m *customMaterialRequisitionItemModel) FindByRequisitionId(ctx context.Context, requisitionId int64) ([]*MaterialRequisitionItem, error) {
	var items []*MaterialRequisitionItem
	query := "SELECT * FROM material_requisition_item WHERE requisition_id = ?"
	err := m.QueryRowsNoCacheCtx(ctx, &items, query, requisitionId)
	return items, err
}

package model

import (
	"context"
	"erp/app/production/rpc/internal/types"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BomModel = (*customBomModel)(nil)

type (
	// BomModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBomModel.
	BomModel interface {
		bomModel
		GetBomList(ctx context.Context, params *types.GetBomListParams) ([]*Bom, int64, error)
		CreateWithDetails(ctx context.Context, bom *Bom, items []*BomItem) error
		UpdateWithDetails(ctx context.Context, bom *Bom, items []*BomItem) error
		DeleteWithDetails(ctx context.Context, bomId int64, itemIds []int64) error
	}

	customBomModel struct {
		*defaultBomModel
		conn sqlx.SqlConn
	}
)

// NewBomModel returns a model for the database table.
func NewBomModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) BomModel {
	return &customBomModel{
		defaultBomModel: newBomModel(conn, c, opts...),
		conn:            conn,
	}
}

func (m *customBomModel) GetBomList(ctx context.Context, params *types.GetBomListParams) ([]*Bom, int64, error) {
	// 构造查询条件
	var conditions []string
	var args []interface{}

	if params.ProductId != 0 {
		conditions = append(conditions, "product_id = ?")
		args = append(args, params.ProductId)
	}
	if params.IsActive != 0 {
		conditions = append(conditions, "is_active = ?")
		args = append(args, params.IsActive)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			whereClause += " AND " + conditions[i]
		}
	}

	// 查询总数
	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", m.table, whereClause)
	err := m.QueryRowNoCacheCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	offset := (params.Page - 1) * params.PageSize
	listQuery := fmt.Sprintf("SELECT * FROM %s %s ORDER BY created_at DESC LIMIT ? OFFSET ?", m.table, whereClause)
	args = append(args, params.PageSize, offset)

	var boms []*Bom
	err = m.QueryRowsNoCacheCtx(ctx, &boms, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	return boms, total, nil
}

func (m *customBomModel) CreateWithDetails(ctx context.Context, bom *Bom, items []*BomItem) error {
	err := m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		query1 := fmt.Sprintf("insert into bom (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", bomRowsExpectAutoSet)
		_, err := session.ExecCtx(ctx, query1,
			bom.Id, bom.BomNo, bom.ProductId, bom.ProductName, bom.Version, bom.UnitCost, bom.IsActive, bom.Remark, bom.CreatedBy, bom.UpdatedBy,
		)
		if err != nil {
			return err
		}

		query2 := fmt.Sprintf("insert into bom_item (%s) values (?, ?, ?, ?, ?, ?, ?, ?)", bomItemRowsExpectAutoSet)

		for _, d := range items {
			_, err := session.ExecCtx(ctx, query2,
				d.Id, d.BomId, d.MaterialId, d.MaterialName, d.Quantity, d.Unit, d.ScrapRate, d.Remark,
			)
			if err != nil {
				return err
			}
		}

		return nil
	})
	return err
}

func (m *customBomModel) UpdateWithDetails(ctx context.Context, bom *Bom, items []*BomItem) error {
	err := m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		// 更新BOM主表
		query1 := fmt.Sprintf("update %s set `version` = ?, `is_active` = ?, `remark` = ?, `updated_by` = ? where `id` = ?", m.table)
		_, err := session.ExecCtx(ctx, query1, bom.Version, bom.IsActive, bom.Remark, bom.UpdatedBy, bom.Id)
		if err != nil {
			return err
		}

		// 删除旧的明细
		query2 := "delete from bom_item where `bom_id` = ?"
		_, err = session.ExecCtx(ctx, query2, bom.Id)
		if err != nil {
			return err
		}

		// 插入新的明细
		if len(items) > 0 {
			query3 := fmt.Sprintf("insert into bom_item (%s) values (?, ?, ?, ?, ?, ?, ?, ?)", bomItemRowsExpectAutoSet)
			for _, item := range items {
				_, err := session.ExecCtx(ctx, query3,
					item.Id, item.BomId, item.MaterialId, item.MaterialName, item.Quantity, item.Unit, item.ScrapRate, item.Remark,
				)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
	return err
}

func (m *customBomModel) DeleteWithDetails(ctx context.Context, bomId int64, itemIds []int64) error {
	err := m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		query1 := fmt.Sprintf("delete from bom where `id` = ?")
		_, err := session.ExecCtx(ctx, query1, bomId)
		if err != nil {
			return err
		}

		query2 := fmt.Sprintf("delete from bom_item where `id` = ?")
		for _, id := range itemIds {
			_, err := session.ExecCtx(ctx, query2, id)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

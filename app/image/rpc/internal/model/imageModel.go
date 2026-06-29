package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ImageModel = (*customImageModel)(nil)

type (
	// ImageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customImageModel.
	ImageModel interface {
		imageModel
		FindByBiz(ctx context.Context, businessType, businessId int64, page, limit int64) ([]*Image, error)
	}

	customImageModel struct {
		*defaultImageModel
	}
)

// NewImageModel returns a model for the database table.
func NewImageModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ImageModel {
	return &customImageModel{
		defaultImageModel: newImageModel(conn, c, opts...),
	}
}

func (m *customImageModel) FindByBiz(ctx context.Context, businessType, businessId int64, page, limit int64) ([]*Image, error) {
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	query := fmt.Sprintf("select %s from %s where `business_type` = ? and `business_id` = ? order by `image_order` asc, `id` desc limit ? offset ?", imageRows, m.table)
	var resp []*Image
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, businessType, businessId, limit, offset)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []*Image{}, nil
	default:
		return nil, err
	}
}

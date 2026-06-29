package productionlogic

import (
	"context"
	"database/sql"
	"fmt"

	"erp/app/production/rpc/internal/svc"
	"erp/app/production/rpc/internal/types"
	"erp/app/production/rpc/production"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateBomItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateBomItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBomItemLogic {
	return &UpdateBomItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateBomItemLogic) UpdateBomItem(in *production.UpdateBomItemReq) (*production.UpdateBomItemResp, error) {
	// 查找现有记录
	bomItem, err := l.svcCtx.BomItemModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	bomItem.MaterialId = in.MaterialId
	bomItem.Quantity = in.Quantity
	bomItem.Unit = sql.NullString{String: in.Unit, Valid: in.Unit != ""}
	bomItem.ScrapRate = in.ScrapRate
	bomItem.Remark = sql.NullString{String: in.Remark, Valid: in.Remark != ""}

	// 更新数据库
	err = l.svcCtx.BomItemModel.Update(l.ctx, bomItem)
	if err != nil {
		return nil, err
	}

	// 清除缓存
	cacheKey := fmt.Sprintf("%s%v", types.CacheErpProductionBomItemIdPrefix, in.Id)
	l.svcCtx.BizRedis.DelCtx(l.ctx, cacheKey)

	return &production.UpdateBomItemResp{}, nil
}

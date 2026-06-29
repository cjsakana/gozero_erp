package productionlogic

import (
	"context"
	"fmt"

	"erp/app/production/rpc/internal/svc"
	"erp/app/production/rpc/internal/types"
	"erp/app/production/rpc/production"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteBomItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteBomItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBomItemLogic {
	return &DeleteBomItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteBomItemLogic) DeleteBomItem(in *production.IdReq) (*production.EmptyResp, error) {
	// 删除 BOM 明细
	err := l.svcCtx.BomItemModel.Delete(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// 清除缓存
	cacheKey := fmt.Sprintf("%s%v", types.CacheErpProductionBomItemIdPrefix, in.Id)
	l.svcCtx.BizRedis.DelCtx(l.ctx, cacheKey)

	return &production.EmptyResp{}, nil
}

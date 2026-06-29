package productionlogic

import (
	"context"
	"fmt"

	"erp/app/production/rpc/internal/svc"
	"erp/app/production/rpc/internal/types"
	"erp/app/production/rpc/production"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteRequisitionItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteRequisitionItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRequisitionItemLogic {
	return &DeleteRequisitionItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteRequisitionItemLogic) DeleteRequisitionItem(in *production.IdReq) (*production.EmptyResp, error) {
	// 删除领料单明细
	err := l.svcCtx.MaterialRequisitionItemModel.Delete(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// 清除缓存
	cacheKey := fmt.Sprintf("%s%v", types.CacheErpProductionMaterialRequisitionItemIdPrefix, in.Id)
	l.svcCtx.BizRedis.DelCtx(l.ctx, cacheKey)

	return &production.EmptyResp{}, nil
}

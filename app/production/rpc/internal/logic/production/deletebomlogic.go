package productionlogic

import (
	"context"
	"erp/app/production/rpc/internal/types"
	"fmt"

	"erp/app/production/rpc/internal/svc"
	"erp/app/production/rpc/production"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteBomLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteBomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBomLogic {
	return &DeleteBomLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteBomLogic) DeleteBom(in *production.IdReq) (*production.EmptyResp, error) {
	var key []string

	// 删除BOM明细
	items, _ := l.svcCtx.BomItemModel.FindByBomId(l.ctx, in.Id)
	ids := make([]int64, len(items))
	for _, item := range items {
		ids = append(ids, item.Id)
		erpProductionBomItemIdKey := fmt.Sprintf("%s%v", types.CacheErpProductionBomItemIdPrefix, item.Id)
		key = append(key, erpProductionBomItemIdKey)
	}

	// 删除BOM主表
	err := l.svcCtx.BomModel.Delete(l.ctx, in.Id)

	err = l.svcCtx.BomModel.DeleteWithDetails(l.ctx, in.Id, ids)
	if err != nil {
		return nil, err
	}

	erpProductionBomIdKey := fmt.Sprintf("%s%v", types.CacheErpProductionBomIdPrefix, in.Id)
	key = append(key, erpProductionBomIdKey)

	l.svcCtx.BizRedis.DelCtx(l.ctx, key...)

	return &production.EmptyResp{}, nil
}

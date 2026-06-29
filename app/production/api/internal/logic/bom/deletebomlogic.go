package bom

import (
	"context"
	"erp/app/production/rpc/production"
	"erp/common/util"

	"erp/app/production/api/internal/svc"
	"erp/app/production/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteBomLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除BOM
func NewDeleteBomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBomLogic {
	return &DeleteBomLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteBomLogic) DeleteBom(req *types.IdReq) (resp *types.DeleteBomResp, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.ProductionRPC.DeleteBom(l.ctx, &production.IdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.DeleteBomResp{}
	return
}

package bomitem

import (
	"context"

	"erp/app/production/api/internal/svc"
	"erp/app/production/api/internal/types"
	"erp/app/production/rpc/production"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteBomItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除BOM明细
func NewDeleteBomItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBomItemLogic {
	return &DeleteBomItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteBomItemLogic) DeleteBomItem(req *types.IdReq) (resp *types.DeleteBomItemResp, err error) {
	// 转换 ID
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}

	// 调用 RPC 删除 BOM 明细
	_, err = l.svcCtx.ProductionRPC.DeleteBomItem(l.ctx, &production.IdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.DeleteBomItemResp{}
	return
}

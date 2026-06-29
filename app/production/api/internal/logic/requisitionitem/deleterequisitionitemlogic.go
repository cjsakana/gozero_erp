package requisitionitem

import (
	"context"

	"erp/app/production/api/internal/svc"
	"erp/app/production/api/internal/types"
	"erp/app/production/rpc/production"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteRequisitionItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除领料单明细
func NewDeleteRequisitionItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRequisitionItemLogic {
	return &DeleteRequisitionItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteRequisitionItemLogic) DeleteRequisitionItem(req *types.IdReq) (resp *types.DeleteRequisitionItemResp, err error) {
	// 转换 ID
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}

	// 调用 RPC 删除领料单明细
	_, err = l.svcCtx.ProductionRPC.DeleteRequisitionItem(l.ctx, &production.IdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.DeleteRequisitionItemResp{}
	return
}

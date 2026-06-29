package inventorylogic

import (
	"context"

	"erp/app/inventory/rpc/internal/svc"
	"erp/app/inventory/rpc/pb"

	"erp/app/inventory/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type DelInventoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelInventoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelInventoryLogic {
	return &DelInventoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelInventoryLogic) DelInventory(in *pb.DelInventoryReq) (*pb.DelInventoryResp, error) {
	err := l.svcCtx.InventoryModel.Delete(l.ctx, in.Id)
	if err != nil {

		return nil, code.AdjustInventoryFail

	}
	return &pb.DelInventoryResp{}, nil
}

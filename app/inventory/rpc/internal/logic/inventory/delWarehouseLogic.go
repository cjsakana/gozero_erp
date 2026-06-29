package inventorylogic

import (
	"context"

	"erp/app/inventory/rpc/internal/svc"
	"erp/app/inventory/rpc/pb"

	"erp/app/inventory/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type DelWarehouseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelWarehouseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelWarehouseLogic {
	return &DelWarehouseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelWarehouseLogic) DelWarehouse(in *pb.DelWarehouseReq) (*pb.DelWarehouseResp, error) {
	err := l.svcCtx.WarehouseModel.Delete(l.ctx, in.Id)
	if err != nil {
		return nil, code.DeleteWarehouseFail
	}

	return &pb.DelWarehouseResp{}, nil
}

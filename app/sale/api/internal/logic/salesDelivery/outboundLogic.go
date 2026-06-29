package salesDelivery

import (
	"context"
	"erp/app/sale/rpc/pb"
	"erp/common/util"
	"erp/common/xtypes"

	"erp/app/sale/api/internal/svc"
	"erp/app/sale/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OutboundLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOutboundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OutboundLogic {
	return &OutboundLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OutboundLogic) Outbound(req *types.OutboundReq) (resp *types.OutboundResp, err error) {
	createdBy, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}
	items := make([]*pb.OutboundDetailItem, 0, len(req.Items))
	for _, v := range req.Items {
		id, err := util.StringToInt64(v.Id)
		if err != nil {
			return nil, err
		}
		batchId, err := util.StringToInt64(v.BatchId)
		if err != nil {
			return nil, err
		}
		items = append(items, &pb.OutboundDetailItem{
			Id:       id,
			Quantity: v.Quantity,
			BatchId:  batchId,
		})
	}
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.SaleRPC.Outbound(l.ctx, &pb.OutboundReq{
		Id:        id,
		CreatedBy: createdBy,
		Items:     items,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.OutboundResp{}
	return
}

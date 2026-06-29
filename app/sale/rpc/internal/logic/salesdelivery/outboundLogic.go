package salesdeliverylogic

import (
	"context"

	"erp/app/sale/rpc/internal/code"
	"erp/app/sale/rpc/internal/svc"
	"erp/app/sale/rpc/internal/types"
	"erp/app/sale/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OutboundLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOutboundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OutboundLogic {
	return &OutboundLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OutboundLogic) Outbound(in *pb.OutboundReq) (*pb.OutboundResp, error) {
	param := &types.OutboundParam{
		Id:        in.Id,
		CreatedBy: in.CreatedBy,
	}

	for _, item := range in.Items {
		param.Items = append(param.Items, types.OutboundDetailParam{
			Id:       item.Id,
			Quantity: item.Quantity,
			BatchId:  item.BatchId,
		})
	}

	keys, err := l.svcCtx.SalesDeliveryModel.Outbound(l.ctx, param)
	if err != nil {
		return nil, code.OutboundFail
	}

	go func(keys []string) {
		// 删掉 缓存
		for range 3 {
			_, err := l.svcCtx.BizRedis.DelCtx(l.ctx, keys...)
			if err == nil {
				return
			}
		}
	}(keys)

	return &pb.OutboundResp{}, nil
}

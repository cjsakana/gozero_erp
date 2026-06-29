package salesOrder

import (
	"context"
	"erp/app/sale/rpc/pb"
	"erp/common/util"

	"erp/app/sale/api/internal/svc"
	"erp/app/sale/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSalesOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSalesOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSalesOrderLogic {
	return &UpdateSalesOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSalesOrderLogic) UpdateSalesOrder(req *types.UpdateSalesOrderRequest) (resp *types.UpdateSalesOrderResponse, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.SaleRPC.UpdateSalesOrder(l.ctx, &pb.UpdateSalesOrderReq{
		Id:          id,
		Status:      req.Status,
		ContractUrl: req.ContractURL,
	})
	if err != nil {
		return nil, err
	}

	return
}

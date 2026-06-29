package salesorderlogic

import (
	"context"
	"database/sql"
	"erp/app/sale/rpc/internal/model"

	"erp/app/sale/rpc/internal/svc"
	"erp/app/sale/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSalesOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateSalesOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSalesOrderLogic {
	return &UpdateSalesOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// rpc UpdateSalesOrderStatus(UpdateSalesOrderStatusReq) returns (UpdateSalesOrderStatusResp);
func (l *UpdateSalesOrderLogic) UpdateSalesOrder(in *pb.UpdateSalesOrderReq) (*pb.UpdateSalesOrderResp, error) {
	if err := l.svcCtx.SalesOrderModel.XUpdate(l.ctx, &model.SalesOrder{
		Id:          in.Id,
		Status:      in.Status,
		ContractUrl: sql.NullString{String: in.ContractUrl, Valid: true},
	}); err != nil {
		return nil, err
	}

	return &pb.UpdateSalesOrderResp{}, nil
}

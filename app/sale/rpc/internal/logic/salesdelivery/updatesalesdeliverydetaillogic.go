package salesdeliverylogic

import (
	"context"
	"database/sql"

	"erp/app/sale/rpc/internal/model"
	"erp/app/sale/rpc/internal/svc"
	"erp/app/sale/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSalesDeliveryDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateSalesDeliveryDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSalesDeliveryDetailLogic {
	return &UpdateSalesDeliveryDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateSalesDeliveryDetailLogic) UpdateSalesDeliveryDetail(in *pb.UpdateSalesDeliveryDetailReq) (*pb.UpdateSalesDeliveryDetailResp, error) {
	// 构造更新数据
	detail := &model.SalesDeliveryDetail{
		Id:          in.Id,
		DeliveryId:  in.DeliveryId,
		ProductId:   in.ProductId,
		ProductName: sql.NullString{String: in.ProductName, Valid: in.ProductName != ""},
		Unit:        in.Unit,
		Quantity:    in.Quantity,
		UnitPrice:   in.UnitPrice,
		Amount:      in.Amount,
		BatchId:     sql.NullInt64{Int64: in.BatchId, Valid: in.BatchId > 0},
	}

	// 更新明细
	err := l.svcCtx.SalesDeliveryDetailModel.XUpdate(l.ctx, detail)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateSalesDeliveryDetailResp{}, nil
}

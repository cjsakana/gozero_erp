package paymentRecord

import (
	"context"
	"erp/app/finance/api/internal/svc"
	"erp/app/finance/api/internal/types"
	"erp/app/finance/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePaymentRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePaymentRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePaymentRecordLogic {
	return &UpdatePaymentRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePaymentRecordLogic) UpdatePaymentRecord(req *types.UpdatePaymentRecordReq) (resp *types.UpdatePaymentRecordResp, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	supplierId, err := util.StringToInt64(req.SupplierId)
	if err != nil {
		return nil, err
	}
	orderId, err := util.StringToInt64(req.OrderId)
	if err != nil {
		return nil, err
	}
	operatorId, err := util.StringToInt64(req.OperatorId)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.FinanceRPC.UpdatePaymentRecord(l.ctx, &pb.UpdatePaymentRecordReq{
		Id:            id,
		SupplierId:    supplierId,
		OrderId:       orderId,
		PaymentType:   req.PaymentType,
		Amount:        req.Amount,
		PaymentDate:   req.PaymentDate,
		PaymentMethod: req.PaymentMethod,
		Status:        req.Status,
		VerifyStatus:  req.VerifyStatus,
		OperatorId:    operatorId,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.UpdatePaymentRecordResp{}
	return
}

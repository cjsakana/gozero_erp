package paymentRecord

import (
	"context"
	"erp/app/finance/api/internal/svc"
	"erp/app/finance/api/internal/types"
	"erp/app/finance/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddPaymentRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddPaymentRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddPaymentRecordLogic {
	return &AddPaymentRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddPaymentRecordLogic) AddPaymentRecord(req *types.AddPaymentRecordReq) (resp *types.AddPaymentRecordResp, err error) {
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
	ret, err := l.svcCtx.FinanceRPC.AddPaymentRecord(l.ctx, &pb.AddPaymentRecordReq{
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

	resp = &types.AddPaymentRecordResp{
		Id: util.Int64ToString(ret.Id),
	}
	return
}

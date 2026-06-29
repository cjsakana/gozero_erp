package receiptRecord

import (
	"context"
	"erp/app/finance/api/internal/svc"
	"erp/app/finance/api/internal/types"
	"erp/app/finance/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddReceiptRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddReceiptRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddReceiptRecordLogic {
	return &AddReceiptRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddReceiptRecordLogic) AddReceiptRecord(req *types.AddReceiptRecordReq) (resp *types.AddReceiptRecordResp, err error) {
	customerId, err := util.StringToInt64(req.CustomerId)
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
	ret, err := l.svcCtx.FinanceRPC.AddReceiptRecord(l.ctx, &pb.AddReceiptRecordReq{
		CustomerId:    customerId,
		OrderId:       orderId,
		ReceiptType:   req.ReceiptType,
		Amount:        req.Amount,
		ReceiptDate:   req.ReceiptDate,
		ReceiptMethod: req.ReceiptMethod,
		Status:        req.Status,
		VerifyStatus:  req.VerifyStatus,
		OperatorId:    operatorId,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.AddReceiptRecordResp{
		Id: util.Int64ToString(ret.Id),
	}
	return
}

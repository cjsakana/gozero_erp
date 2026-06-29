package receiptRecord

import (
	"context"
	"erp/app/finance/api/internal/svc"
	"erp/app/finance/api/internal/types"
	"erp/app/finance/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateReceiptRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateReceiptRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateReceiptRecordLogic {
	return &UpdateReceiptRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateReceiptRecordLogic) UpdateReceiptRecord(req *types.UpdateReceiptRecordReq) (resp *types.UpdateReceiptRecordResp, err error) {
	id, err := util.StringToInt64(req.Id)
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
	_, err = l.svcCtx.FinanceRPC.UpdateReceiptRecord(l.ctx, &pb.UpdateReceiptRecordReq{
		Id:            id,
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

	resp = &types.UpdateReceiptRecordResp{}
	return
}

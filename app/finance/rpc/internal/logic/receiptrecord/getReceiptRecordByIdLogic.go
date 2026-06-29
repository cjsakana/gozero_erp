package receiptrecordlogic

import (
	"context"

	"erp/app/finance/rpc/internal/code"
	"erp/app/finance/rpc/internal/model"
	"erp/app/finance/rpc/internal/svc"
	"erp/app/finance/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetReceiptRecordByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	Logger logx.Logger
}

func NewGetReceiptRecordByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetReceiptRecordByIdLogic {
	return &GetReceiptRecordByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetReceiptRecordByIdLogic) GetReceiptRecordById(in *pb.GetReceiptRecordByIdReq) (*pb.GetReceiptRecordByIdResp, error) {
	rr, err := l.svcCtx.ReceiptRecordModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, code.GetReceiptRecordFail
		}
		return nil, code.GetReceiptRecordFail
	}

	return &pb.GetReceiptRecordByIdResp{
		ReceiptRecord: &pb.ReceiptRecord{
			Id:            rr.Id,
			ReceiptNo:     rr.ReceiptNo,
			CustomerId:    rr.CustomerId,
			OrderId:       rr.OrderId.Int64,
			ReceiptType:   rr.ReceiptType,
			Amount:        rr.Amount,
			ReceiptDate:   rr.ReceiptDate.Unix(),
			ReceiptMethod: rr.ReceiptMethod.String,
			Status:        rr.Status,
			VerifyStatus:  rr.VerifyStatus,
			OperatorId:    rr.OperatorId,
			CreatedAt:     rr.CreatedAt.Unix(),
		},
	}, nil
}

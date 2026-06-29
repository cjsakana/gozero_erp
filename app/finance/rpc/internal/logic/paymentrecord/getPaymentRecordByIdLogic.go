package paymentrecordlogic

import (
	"context"

	"erp/app/finance/rpc/internal/svc"
	"erp/app/finance/rpc/pb"

	"erp/app/finance/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

type GetPaymentRecordByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPaymentRecordByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaymentRecordByIdLogic {
	return &GetPaymentRecordByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPaymentRecordByIdLogic) GetPaymentRecordById(in *pb.GetPaymentRecordByIdReq) (*pb.GetPaymentRecordByIdResp, error) {
	paymentRecord, err := l.svcCtx.PaymentRecordModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.PaymentRecordNotFound
		}
		return nil, code.PaymentRecordNotFound
	}

	resp := &pb.GetPaymentRecordByIdResp{
		PaymentRecord: &pb.PaymentRecord{
			Id:            paymentRecord.Id,
			PaymentNo:     paymentRecord.PaymentNo,
			SupplierId:    paymentRecord.SupplierId,
			OrderId:       paymentRecord.OrderId.Int64,
			PaymentType:   paymentRecord.PaymentType,
			Amount:        paymentRecord.Amount,
			PaymentDate:   paymentRecord.PaymentDate.Unix(),
			PaymentMethod: paymentRecord.PaymentMethod.String,
			Status:        paymentRecord.Status,
			VerifyStatus:  paymentRecord.VerifyStatus,
			OperatorId:    paymentRecord.OperatorId,
			CreatedAt:     paymentRecord.CreatedAt.Unix(),
		},
	}
	return resp, nil
}

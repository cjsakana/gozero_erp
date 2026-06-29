package paymentrecordlogic

import (
	"context"
	"database/sql"
	"time"

	"erp/app/finance/rpc/internal/model"
	"erp/app/finance/rpc/internal/svc"
	"erp/app/finance/rpc/pb"

	"erp/app/finance/rpc/internal/code"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePaymentRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePaymentRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePaymentRecordLogic {
	return &UpdatePaymentRecordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdatePaymentRecordLogic) UpdatePaymentRecord(in *pb.UpdatePaymentRecordReq) (*pb.UpdatePaymentRecordResp, error) {
	data := &model.PaymentRecord{
		Id:            in.Id,
		SupplierId:    in.SupplierId,
		OrderId:       sql.NullInt64{Int64: in.OrderId, Valid: in.OrderId != 0},
		PaymentType:   in.PaymentType,
		Amount:        in.Amount,
		PaymentDate:   time.Unix(in.PaymentDate, 0),
		PaymentMethod: sql.NullString{String: in.PaymentMethod, Valid: true},
		Status:        in.Status,
		VerifyStatus:  in.VerifyStatus,
		OperatorId:    in.OperatorId,
	}

	err := l.svcCtx.PaymentRecordModel.Update(l.ctx, data)
	if err != nil {

		return nil, code.UpdatePaymentRecordFail

	}

	return &pb.UpdatePaymentRecordResp{}, nil
}

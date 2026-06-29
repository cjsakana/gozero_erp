package paymentrecordlogic

import (
	"context"
	"database/sql"
	"erp/common/util"
	"time"

	"erp/app/finance/rpc/internal/model"
	"erp/app/finance/rpc/internal/svc"
	"erp/app/finance/rpc/pb"

	"erp/app/finance/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddPaymentRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddPaymentRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddPaymentRecordLogic {
	return &AddPaymentRecordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------paymentRecord-----------------------
func (l *AddPaymentRecordLogic) AddPaymentRecord(in *pb.AddPaymentRecordReq) (*pb.AddPaymentRecordResp, error) {
	id := util.GenerateSnowflake()
	paymentNo := util.GenerateNo("PAY")
	data := &model.PaymentRecord{
		Id:            id,
		PaymentNo:     paymentNo,
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

	_, err := l.svcCtx.PaymentRecordModel.Insert(l.ctx, data)
	if err != nil {
		return nil, code.AddPaymentRecordFail
	}

	return &pb.AddPaymentRecordResp{Id: id}, nil
}

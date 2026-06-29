package salarypaymentlogic

import (
	"context"
	"database/sql"
	"time"

	"erp/app/finance/rpc/internal/code"
	"erp/app/finance/rpc/internal/model"
	"erp/app/finance/rpc/internal/svc"
	"erp/app/finance/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSalaryPaymentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	Logger logx.Logger
}

func NewUpdateSalaryPaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSalaryPaymentLogic {
	return &UpdateSalaryPaymentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateSalaryPaymentLogic) UpdateSalaryPayment(in *pb.UpdateSalaryPaymentReq) (*pb.UpdateSalaryPaymentResp, error) {
	data := &model.SalaryPayment{
		Id:            in.Id,
		Bonus:         in.Bonus,
		OvertimePay:   in.OvertimePay,
		Deduction:     in.Deduction,
		NetPayment:    in.NetPayment,
		PaymentDate:   time.Unix(in.PaymentDate, 0),
		PaymentMethod: sql.NullString{String: in.PaymentMethod, Valid: in.PaymentMethod != ""},
		BankAccount:   sql.NullString{String: in.BankAccount, Valid: in.BankAccount != ""},
		ReferenceNo:   sql.NullString{String: in.ReferenceNo, Valid: in.ReferenceNo != ""},
		Status:        in.Status,
	}

	if err := l.svcCtx.SalaryPaymentModel.Update(l.ctx, data); err != nil {
		return nil, code.UpdateSalaryPaymentFail
	}
	return &pb.UpdateSalaryPaymentResp{}, nil
}

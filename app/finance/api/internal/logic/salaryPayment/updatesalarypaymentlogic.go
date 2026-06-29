package salaryPayment

import (
	"context"
	"erp/app/finance/api/internal/svc"
	"erp/app/finance/api/internal/types"
	"erp/app/finance/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSalaryPaymentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSalaryPaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSalaryPaymentLogic {
	return &UpdateSalaryPaymentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSalaryPaymentLogic) UpdateSalaryPayment(req *types.UpdateSalaryPaymentReq) (resp *types.UpdateSalaryPaymentResp, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.FinanceRPC.UpdateSalaryPayment(l.ctx, &pb.UpdateSalaryPaymentReq{
		Id:            id,
		Bonus:         req.Bonus,
		OvertimePay:   req.OvertimePay,
		Deduction:     req.Deduction,
		NetPayment:    req.NetPayment,
		PaymentDate:   req.PaymentDate,
		PaymentMethod: req.PaymentMethod,
		BankAccount:   req.BankAccount,
		ReferenceNo:   req.ReferenceNo,
		Status:        req.Status,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.UpdateSalaryPaymentResp{}
	return
}

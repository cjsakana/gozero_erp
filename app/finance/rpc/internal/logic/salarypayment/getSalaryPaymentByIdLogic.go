package salarypaymentlogic

import (
	"context"

	"erp/app/finance/rpc/internal/code"
	"erp/app/finance/rpc/internal/model"
	"erp/app/finance/rpc/internal/svc"
	"erp/app/finance/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSalaryPaymentByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	Logger logx.Logger
}

func NewGetSalaryPaymentByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSalaryPaymentByIdLogic {
	return &GetSalaryPaymentByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSalaryPaymentByIdLogic) GetSalaryPaymentById(in *pb.GetSalaryPaymentByIdReq) (*pb.GetSalaryPaymentByIdResp, error) {
	sp, err := l.svcCtx.SalaryPaymentModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, code.GetSalaryPaymentFail
		}
		return nil, code.GetSalaryPaymentFail
	}

	return &pb.GetSalaryPaymentByIdResp{
		SalaryPayment: &pb.SalaryPayment{
			Id:            sp.Id,
			PayrollId:     sp.PayrollId.Int64,
			EmployeeId:    sp.EmployeeId,
			EmployeeName:  sp.EmployeeName,
			DepartmentId:  sp.DepartmentId.Int64,
			PaymentMonth:  sp.PaymentMonth.Unix(),
			BaseSalary:    sp.BaseSalary.Float64,
			Bonus:         sp.Bonus,
			OvertimePay:   sp.OvertimePay,
			Deduction:     sp.Deduction,
			NetPayment:    sp.NetPayment,
			PaymentDate:   sp.PaymentDate.Unix(),
			PaymentMethod: sp.PaymentMethod.String,
			BankAccount:   sp.BankAccount.String,
			ReferenceNo:   sp.ReferenceNo.String,
			Status:        sp.Status,
			CreatedAt:     sp.CreatedAt.Unix(),
			UpdatedAt:     sp.UpdatedAt.Unix(),
		},
	}, nil
}

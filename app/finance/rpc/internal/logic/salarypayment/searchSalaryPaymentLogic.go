package salarypaymentlogic

import (
	"context"

	"erp/app/finance/rpc/internal/model"
	"erp/app/finance/rpc/internal/svc"
	"erp/app/finance/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchSalaryPaymentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	Logger logx.Logger
}

func NewSearchSalaryPaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchSalaryPaymentLogic {
	return &SearchSalaryPaymentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchSalaryPaymentLogic) SearchSalaryPayment(in *pb.SearchSalaryPaymentReq) (*pb.SearchSalaryPaymentResp, error) {
	list, total, err := l.svcCtx.SalaryPaymentModel.Search(l.ctx,
		in.EmployeeId, in.DepartmentId, in.Status, in.PaymentMonth, in.Page, in.Limit,
	)
	if err != nil {
		return nil, model.ErrNotFound
	}

	items := make([]*pb.SalaryPayment, 0, len(list))
	for _, sp := range list {
		items = append(items, &pb.SalaryPayment{
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
		})
	}
	return &pb.SearchSalaryPaymentResp{SalaryPayment: items, Total: total}, nil
}

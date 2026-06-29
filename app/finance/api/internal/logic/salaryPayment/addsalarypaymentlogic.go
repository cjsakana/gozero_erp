package salaryPayment

import (
	"context"
	"erp/app/finance/api/internal/svc"
	"erp/app/finance/api/internal/types"
	"erp/app/finance/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddSalaryPaymentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddSalaryPaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSalaryPaymentLogic {
	return &AddSalaryPaymentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddSalaryPaymentLogic) AddSalaryPayment(req *types.AddSalaryPaymentReq) (resp *types.AddSalaryPaymentResp, err error) {
	payrollId, err := util.StringToInt64(req.PayrollId)
	if err != nil {
		return nil, err
	}
	employeeId, err := util.StringToInt64(req.EmployeeId)
	if err != nil {
		return nil, err
	}
	departmentId, err := util.StringToInt64(req.DepartmentId)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.FinanceRPC.AddSalaryPayment(l.ctx, &pb.AddSalaryPaymentReq{
		PayrollId:     payrollId,
		EmployeeId:    employeeId,
		EmployeeName:  req.EmployeeName,
		DepartmentId:  departmentId,
		PaymentMonth:  req.PaymentMonth,
		BaseSalary:    req.BaseSalary,
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

	resp = &types.AddSalaryPaymentResp{
		Id: util.Int64ToString(ret.Id),
	}
	return
}

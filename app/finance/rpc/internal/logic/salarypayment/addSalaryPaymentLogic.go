package salarypaymentlogic

import (
	"context"
	"database/sql"
	"erp/common/util"
	"time"

	"erp/app/finance/rpc/internal/code"
	"erp/app/finance/rpc/internal/model"
	"erp/app/finance/rpc/internal/svc"
	"erp/app/finance/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddSalaryPaymentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	Logger logx.Logger
}

func NewAddSalaryPaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSalaryPaymentLogic {
	return &AddSalaryPaymentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------salaryPayment-----------------------
func (l *AddSalaryPaymentLogic) AddSalaryPayment(in *pb.AddSalaryPaymentReq) (*pb.AddSalaryPaymentResp, error) {
	id := util.GenerateSnowflake()
	data := &model.SalaryPayment{
		Id:            id,
		PayrollId:     sql.NullInt64{Int64: in.PayrollId, Valid: in.PayrollId != 0},
		EmployeeId:    in.EmployeeId,
		EmployeeName:  in.EmployeeName,
		DepartmentId:  sql.NullInt64{Int64: in.DepartmentId, Valid: in.DepartmentId != 0},
		PaymentMonth:  time.Unix(in.PaymentMonth, 0),
		BaseSalary:    sql.NullFloat64{Float64: in.BaseSalary, Valid: in.BaseSalary != 0},
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

	_, err := l.svcCtx.SalaryPaymentModel.Insert(l.ctx, data)
	if err != nil {
		return nil, code.AddSalaryPaymentFail
	}

	return &pb.AddSalaryPaymentResp{Id: id}, nil
}

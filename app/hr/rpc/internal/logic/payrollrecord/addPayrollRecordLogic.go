package payrollrecordlogic

import (
	"context"
	"database/sql"
	"erp/app/hr/rpc/internal/model"
	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"
	"erp/common/util"
	"time"

	"erp/app/hr/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddPayrollRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddPayrollRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddPayrollRecordLogic {
	return &AddPayrollRecordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------payrollRecord-----------------------
func (l *AddPayrollRecordLogic) AddPayrollRecord(in *pb.AddPayrollRecordReq) (*pb.AddPayrollRecordResp, error) {
	// 根据员工ID查询员工信息
	one, err := l.svcCtx.EmployeeDetailModel.FindOne(l.ctx, in.EmployeeId)
	if err != nil {
		return nil, code.AddPayrollFail
	}

	netSalary := one.Salary.Float64 + in.Bonus - in.Deductions

	// 生成雪花ID
	id := util.GenerateSnowflake()
	_, err = l.svcCtx.PayrollRecordModel.Insert(l.ctx, &model.PayrollRecord{
		Id:           id,
		EmployeeId:   in.EmployeeId, // 使用员工ID（新版主键）
		PaymentMonth: time.Unix(in.PaymentMonth, 0),
		BaseSalary:   one.Salary,
		Bonus:        in.Bonus,
		Deductions:   in.Deductions,
		NetSalary:    sql.NullFloat64{Float64: netSalary, Valid: true},
		CalculatedBy: sql.NullInt64{Int64: in.CalculatedBy, Valid: in.CalculatedBy != 0},
		CalculatedAt: sql.NullTime{Time: time.Unix(in.CalculatedAt, 0), Valid: in.CalculatedAt != 0},
		Status:       1, // 审批中
		Description:  sql.NullString{String: in.Description, Valid: in.Description != ""},
		PaymentAt:    sql.NullTime{Valid: false},
	})
	if err != nil {
		return nil, code.AddPayrollFail
	}
	return &pb.AddPayrollRecordResp{
		Id: id,
	}, nil
}

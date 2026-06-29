package payrollrecordlogic

import (
	"context"
	"database/sql"
	"erp/app/hr/rpc/internal/code"
	"erp/app/hr/rpc/internal/model"
	"time"

	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

type UpdatePayrollRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePayrollRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePayrollRecordLogic {
	return &UpdatePayrollRecordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdatePayrollRecordLogic) UpdatePayrollRecord(in *pb.UpdatePayrollRecordReq) (*pb.UpdatePayrollRecordResp, error) {
	payrollRecord, err := l.svcCtx.PayrollRecordModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.PayrollNotFound
		}
		return nil, code.GetPayrollFail
	}
	one, err := l.svcCtx.EmployeeDetailModel.FindOne(l.ctx, payrollRecord.EmployeeId)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.EmployeeNotFound
		}
		return nil, code.GetEmployeeFail
	}
	netSalary := one.Salary.Float64 + in.Bonus - in.Deductions
	err = l.svcCtx.PayrollRecordModel.XUpdate(l.ctx, &model.PayrollRecord{
		Id:           in.Id,
		Bonus:        in.Bonus,
		Deductions:   in.Deductions,
		NetSalary:    sql.NullFloat64{Float64: netSalary, Valid: true},
		CalculatedBy: sql.NullInt64{Int64: in.CalculatedBy, Valid: true},
		CalculatedAt: sql.NullTime{Time: time.Unix(in.CalculatedAt, 0), Valid: true},
		Status:       in.Status,
		Description:  sql.NullString{String: in.Description, Valid: true},
		PaymentAt:    sql.NullTime{Time: time.Unix(in.PaymentAt, 0), Valid: true},
	})
	if err != nil {
		return nil, code.ApprovePayrollFail
	}

	return &pb.UpdatePayrollRecordResp{}, nil
}

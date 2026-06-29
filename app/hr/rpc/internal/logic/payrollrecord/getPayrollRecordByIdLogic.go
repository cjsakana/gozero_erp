package payrollrecordlogic

import (
	"context"
	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"

	"erp/app/hr/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

type GetPayrollRecordByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPayrollRecordByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPayrollRecordByIdLogic {
	return &GetPayrollRecordByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPayrollRecordByIdLogic) GetPayrollRecordById(in *pb.GetPayrollRecordByIdReq) (*pb.GetPayrollRecordByIdResp, error) {
	one, err := l.svcCtx.PayrollRecordModel.FindOne(l.ctx, in.Id)

	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.PayrollNotFound
		}
		return nil, code.PayrollNotFound
	}

	return &pb.GetPayrollRecordByIdResp{
		PayrollRecord: &pb.PayrollRecord{
			Id:           one.Id,
			EmployeeId:   one.EmployeeId,
			BaseSalary:   one.BaseSalary.Float64,
			Bonus:        one.Bonus,
			Deductions:   one.Deductions,
			NetSalary:    one.NetSalary.Float64,
			CalculatedBy: one.CalculatedBy.Int64,
			CalculatedAt: one.CalculatedAt.Time.Unix(),
			Status:       one.Status,
			Description:  one.Description.String,
			PaymentAt:    one.PaymentAt.Time.Unix(),
			CreatedAt:    one.CreatedAt.Unix(),
		},
	}, nil
}

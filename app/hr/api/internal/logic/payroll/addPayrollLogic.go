package payroll

import (
	"context"
	"erp/app/hr/rpc/pb"
	"erp/common/util"
	"erp/common/xtypes"
	"time"

	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddPayrollLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddPayrollLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddPayrollLogic {
	return &AddPayrollLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddPayrollLogic) AddPayroll(req *types.AddPayrollRequest) (resp *types.AddPayrollResponse, err error) {
	calculatedById, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}

	employeeId, err := util.StringToInt64(req.EmployeeId)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.HrRPC.AddPayrollRecord(l.ctx, &pb.AddPayrollRecordReq{
		EmployeeId:   employeeId,
		PaymentMonth: req.PaymentMonth,
		Bonus:        req.Bonus,
		Deductions:   req.Deductions,
		Description:  req.Description,
		CalculatedBy: calculatedById,
		CalculatedAt: time.Now().Unix(),
	})
	if err != nil {
		return nil, err
	}

	resp = &types.AddPayrollResponse{
		Id: util.Int64ToString(ret.Id),
	}
	return
}

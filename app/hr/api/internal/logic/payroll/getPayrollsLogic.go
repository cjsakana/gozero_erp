package payroll

import (
	"context"
	"erp/app/hr/rpc/pb"
	"erp/common/util"

	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPayrollsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPayrollsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPayrollsLogic {
	return &GetPayrollsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPayrollsLogic) GetPayrolls(req *types.GetPayrollsRequest) (resp *types.GetPayrollsResponse, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.HrRPC.GetPayrollRecordById(l.ctx, &pb.GetPayrollRecordByIdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	// 查询员工信息
	employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
		Id: ret.PayrollRecord.EmployeeId,
	})
	if err != nil {
		return nil, err
	}

	// 查询核算人信息（如果有）
	var calculatedByNo, calculatedByName string
	if ret.PayrollRecord.CalculatedBy > 0 {
		calculatorDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
			Id: ret.PayrollRecord.CalculatedBy,
		})
		if err != nil {
			logx.Errorf("查询核算人信息失败: calculatedBy=%d, err=%v", ret.PayrollRecord.CalculatedBy, err)
		} else {
			calculatedByNo = calculatorDetail.EmployeeNonSensitiveDetail.EmployeeNo
			calculatedByName = calculatorDetail.EmployeeNonSensitiveDetail.Name
		}
	}

	resp = &types.GetPayrollsResponse{
		PayrollRecord: types.PayrollRecord{
			Id:               util.Int64ToString(ret.PayrollRecord.Id),
			EmployeeId:       util.Int64ToString(ret.PayrollRecord.EmployeeId),
			EmployeeNo:       employeeDetail.EmployeeNonSensitiveDetail.EmployeeNo,
			EmployeeName:     employeeDetail.EmployeeNonSensitiveDetail.Name,
			PaymentMonth:     ret.PayrollRecord.PaymentMonth,
			BaseSalary:       ret.PayrollRecord.BaseSalary,
			Bonus:            ret.PayrollRecord.Bonus,
			Deductions:       ret.PayrollRecord.Deductions,
			NetSalary:        ret.PayrollRecord.NetSalary,
			CalculatedById:   util.Int64ToString(ret.PayrollRecord.CalculatedBy),
			CalculatedByNo:   calculatedByNo,
			CalculatedByName: calculatedByName,
			CalculatedAt:     ret.PayrollRecord.CalculatedAt,
			Status:           ret.PayrollRecord.Status,
			Description:      ret.PayrollRecord.Description,
			PaymentAt:        ret.PayrollRecord.PaymentAt,
			CreatedAt:        ret.PayrollRecord.CreatedAt,
		},
	}

	return
}

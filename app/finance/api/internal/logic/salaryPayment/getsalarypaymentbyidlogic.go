package salaryPayment

import (
	"context"
	"erp/app/finance/api/internal/svc"
	"erp/app/finance/api/internal/types"
	"erp/app/finance/rpc/pb"
	hrpb "erp/app/hr/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSalaryPaymentByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSalaryPaymentByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSalaryPaymentByIdLogic {
	return &GetSalaryPaymentByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSalaryPaymentByIdLogic) GetSalaryPaymentById(req *types.GetSalaryPaymentByIdReq) (resp *types.GetSalaryPaymentByIdResp, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.FinanceRPC.GetSalaryPaymentById(l.ctx, &pb.GetSalaryPaymentByIdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	// 获取员工工号
	var employeeNo string
	if ret.SalaryPayment.EmployeeId > 0 {
		empResp, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &hrpb.GetEmployeeDetailByIdReq{
			Id: ret.SalaryPayment.EmployeeId,
		})
		if err != nil {
			logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", ret.SalaryPayment.EmployeeId, err)
		} else if empResp.EmployeeNonSensitiveDetail != nil {
			employeeNo = empResp.EmployeeNonSensitiveDetail.EmployeeNo
		}
	}

	// 获取部门名称
	var departmentName string
	if ret.SalaryPayment.DepartmentId > 0 {
		deptResp, err := l.svcCtx.HrRPC.DepartmentZrpcClient.GetDepartmentById(l.ctx, &hrpb.GetDepartmentByIdReq{
			Id: ret.SalaryPayment.DepartmentId,
		})
		if err != nil {
			logx.Errorf("查询部门信息失败: departmentId=%d, err=%v", ret.SalaryPayment.DepartmentId, err)
		} else if deptResp.Department != nil {
			departmentName = deptResp.Department.Name
		}
	}

	resp = &types.GetSalaryPaymentByIdResp{
		SalaryPayment: types.SalaryPayment{
			Id:             util.Int64ToString(ret.SalaryPayment.Id),
			PayrollId:      util.Int64ToString(ret.SalaryPayment.PayrollId),
			EmployeeId:     util.Int64ToString(ret.SalaryPayment.EmployeeId),
			EmployeeNo:     employeeNo,
			EmployeeName:   ret.SalaryPayment.EmployeeName,
			DepartmentId:   util.Int64ToString(ret.SalaryPayment.DepartmentId),
			DepartmentName: departmentName,
			PaymentMonth:   ret.SalaryPayment.PaymentMonth,
			BaseSalary:     ret.SalaryPayment.BaseSalary,
			Bonus:          ret.SalaryPayment.Bonus,
			OvertimePay:    ret.SalaryPayment.OvertimePay,
			Deduction:      ret.SalaryPayment.Deduction,
			NetPayment:     ret.SalaryPayment.NetPayment,
			PaymentDate:    ret.SalaryPayment.PaymentDate,
			PaymentMethod:  ret.SalaryPayment.PaymentMethod,
			BankAccount:    ret.SalaryPayment.BankAccount,
			ReferenceNo:    ret.SalaryPayment.ReferenceNo,
			Status:         ret.SalaryPayment.Status,
			CreatedAt:      ret.SalaryPayment.CreatedAt,
			UpdatedAt:      ret.SalaryPayment.UpdatedAt,
		},
	}
	return
}

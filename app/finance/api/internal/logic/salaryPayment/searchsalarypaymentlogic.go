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

type SearchSalaryPaymentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchSalaryPaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchSalaryPaymentLogic {
	return &SearchSalaryPaymentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchSalaryPaymentLogic) SearchSalaryPayment(req *types.SearchSalaryPaymentReq) (resp *types.SearchSalaryPaymentResp, err error) {
	employeeId, err := util.StringToInt64(req.EmployeeId)
	if err != nil {
		return nil, err
	}
	departmentId, err := util.StringToInt64(req.DepartmentId)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.FinanceRPC.SearchSalaryPayment(l.ctx, &pb.SearchSalaryPaymentReq{
		Page:         req.Page,
		Limit:        req.Limit,
		EmployeeId:   employeeId,
		DepartmentId: departmentId,
		PaymentMonth: req.PaymentMonth,
		Status:       req.Status,
	})
	if err != nil {
		return nil, err
	}

	// 批量获取员工工号（去重）
	employeeNoMap := make(map[int64]string)
	for _, sp := range ret.SalaryPayment {
		if sp.EmployeeId > 0 {
			if _, ok := employeeNoMap[sp.EmployeeId]; !ok {
				empResp, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &hrpb.GetEmployeeDetailByIdReq{Id: sp.EmployeeId})
				if err != nil {
					logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", sp.EmployeeId, err)
					employeeNoMap[sp.EmployeeId] = ""
				} else if empResp.EmployeeNonSensitiveDetail != nil {
					employeeNoMap[sp.EmployeeId] = empResp.EmployeeNonSensitiveDetail.EmployeeNo
				}
			}
		}
	}

	// 批量获取部门名称（去重）
	departmentMap := make(map[int64]string)
	for _, sp := range ret.SalaryPayment {
		if sp.DepartmentId > 0 {
			if _, ok := departmentMap[sp.DepartmentId]; !ok {
				deptResp, err := l.svcCtx.HrRPC.DepartmentZrpcClient.GetDepartmentById(l.ctx, &hrpb.GetDepartmentByIdReq{Id: sp.DepartmentId})
				if err != nil {
					logx.Errorf("查询部门信息失败: departmentId=%d, err=%v", sp.DepartmentId, err)
					departmentMap[sp.DepartmentId] = ""
				} else if deptResp.Department != nil {
					departmentMap[sp.DepartmentId] = deptResp.Department.Name
				}
			}
		}
	}

	list := make([]*types.SalaryPayment, 0, len(ret.SalaryPayment))
	for _, sp := range ret.SalaryPayment {
		var empNo, deptName string
		if no, ok := employeeNoMap[sp.EmployeeId]; ok {
			empNo = no
		}
		if name, ok := departmentMap[sp.DepartmentId]; ok {
			deptName = name
		}
		item := &types.SalaryPayment{
			Id:             util.Int64ToString(sp.Id),
			PayrollId:      util.Int64ToString(sp.PayrollId),
			EmployeeId:     util.Int64ToString(sp.EmployeeId),
			EmployeeNo:     empNo,
			EmployeeName:   sp.EmployeeName,
			DepartmentId:   util.Int64ToString(sp.DepartmentId),
			DepartmentName: deptName,
			PaymentMonth:   sp.PaymentMonth,
			BaseSalary:     sp.BaseSalary,
			Bonus:          sp.Bonus,
			OvertimePay:    sp.OvertimePay,
			Deduction:      sp.Deduction,
			NetPayment:     sp.NetPayment,
			PaymentDate:    sp.PaymentDate,
			PaymentMethod:  sp.PaymentMethod,
			BankAccount:    sp.BankAccount,
			ReferenceNo:    sp.ReferenceNo,
			Status:         sp.Status,
			CreatedAt:      sp.CreatedAt,
			UpdatedAt:      sp.UpdatedAt,
		}

		list = append(list, item)
	}

	resp = &types.SearchSalaryPaymentResp{
		SalaryPayment: list,
		Total:         ret.Total,
	}
	return
}

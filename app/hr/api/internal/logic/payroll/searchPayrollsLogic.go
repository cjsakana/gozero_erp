package payroll

import (
	"context"
	"erp/app/hr/rpc/client/employeedetail"
	"erp/app/hr/rpc/pb"
	"erp/common/util"

	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchPayrollsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchPayrollsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchPayrollsLogic {
	return &SearchPayrollsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchPayrollsLogic) SearchPayrolls(req *types.SearchPayrollsRequest) (resp *types.SearchPayrollsResponse, err error) {
	employeeId, err := util.StringToInt64(req.EmployeeId)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.HrRPC.SearchPayrollRecord(l.ctx, &pb.SearchPayrollRecordReq{
		Page:                req.Page,
		Limit:               req.Limit,
		PaymentMonth:        req.PaymentMonth,
		EmployeeId:          employeeId,
		Status:              req.Status,
		CalculatedBy:        0, // 暂不支持按核算人查询
		StartCalculatedDate: req.StartCalculatedAt,
		EndCalculatedDate:   req.EndCalculatedAt,
		StartPaymentDate:    req.StartPaymentDate,
		EndPaymentDate:      req.EndPaymentDate,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.SearchPayrollsResponse{
		Total: ret.Total,
	}

	// 批量查询员工和核算人信息
	employeeMap := make(map[int64]*employeedetail.EmployeeNonSensitiveDetail)
	for _, v := range ret.PayrollRecord {
		// 查询员工信息
		if _, ok := employeeMap[v.EmployeeId]; !ok {
			employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
				Id: v.EmployeeId,
			})
			if err != nil {
				logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", v.EmployeeId, err)
				continue
			}
			employeeMap[v.EmployeeId] = employeeDetail.EmployeeNonSensitiveDetail
		}

		// 查询核算人信息（如果有）
		if v.CalculatedBy > 0 {
			if _, ok := employeeMap[v.CalculatedBy]; !ok {
				calculatorDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
					Id: v.CalculatedBy,
				})
				if err != nil {
					logx.Errorf("查询核算人信息失败: calculatedBy=%d, err=%v", v.CalculatedBy, err)
				} else {
					employeeMap[v.CalculatedBy] = calculatorDetail.EmployeeNonSensitiveDetail
				}
			}
		}

		resp.List = append(resp.List, &types.PayrollRecord{
			Id:               util.Int64ToString(v.Id),
			EmployeeId:       util.Int64ToString(v.EmployeeId),
			EmployeeNo:       employeeMap[v.EmployeeId].EmployeeNo,
			EmployeeName:     employeeMap[v.EmployeeId].Name,
			PaymentMonth:     v.PaymentMonth,
			BaseSalary:       v.BaseSalary,
			Bonus:            v.Bonus,
			Deductions:       v.Deductions,
			NetSalary:        v.NetSalary,
			CalculatedById:   util.Int64ToString(v.CalculatedBy),
			CalculatedByNo:   employeeMap[v.CalculatedBy].EmployeeNo,
			CalculatedByName: employeeMap[v.CalculatedBy].Name,
			CalculatedAt:     v.CalculatedAt,
			Status:           v.Status,
			Description:      v.Description,
			PaymentAt:        v.PaymentAt,
			CreatedAt:        v.CreatedAt,
		})
	}

	return
}

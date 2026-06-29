package employeeDetail

import (
	"context"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"erp/app/hr/rpc/pb"
	"erp/common/util"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetDepartmentEmployeesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDepartmentEmployeesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDepartmentEmployeesLogic {
	return &GetDepartmentEmployeesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDepartmentEmployeesLogic) GetDepartmentEmployees(req *types.GetDepartmentEmployeesRequest) (resp *types.GetDepartmentEmployeesResponse, err error) {
	departmentId, err := util.StringToInt64(req.DepartmentId)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.SearchEmployeeDetail(l.ctx, &pb.SearchEmployeeDetailReq{
		Page:         req.Page,
		Limit:        req.Limit,
		DepartmentId: departmentId,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.GetDepartmentEmployeesResponse{
		Total: ret.Total,
	}
	var list []*types.EmployeeDetail
	for _, v := range ret.EmployeeNonSensitiveDetail {
		departmentById, err2 := l.svcCtx.HrRPC.DepartmentZrpcClient.GetDepartmentById(l.ctx, &pb.GetDepartmentByIdReq{
			Id: v.DepartmentId,
		})
		if err2 != nil {
			return nil, err2
		}
		positionById, err2 := l.svcCtx.HrRPC.PositionZrpcClient.GetPositionById(l.ctx, &pb.GetPositionByIdReq{
			Id: v.PositionId,
		})
		if err2 != nil {
			return nil, err2
		}

		list = append(list, &types.EmployeeDetail{
			Employee: types.Employee{
				Id:           util.Int64ToString(v.Id),
				EmployeeNo:   v.EmployeeNo,
				Name:         v.Name,
				Gender:       v.Gender,
				BirthDate:    v.BirthDate,
				DepartmentId: util.Int64ToString(v.DepartmentId),
				PositionId:   util.Int64ToString(v.PositionId),
				Salary:       v.Salary,
				HireDate:     v.HireDate,
				LeaveDate:    v.LeaveDate,
			},
			DepartmentName: departmentById.Department.Name,
			PositionName:   positionById.Position.Name,
		})
	}
	resp.List = list
	return
}

package employeeDetail

import (
	"context"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"erp/app/hr/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchEmployeeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchEmployeeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchEmployeeLogic {
	return &SearchEmployeeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchEmployeeLogic) SearchEmployee(req *types.SearchEmployeeRequest) (resp *types.SearchEmployeeResponse, err error) {
	var departmentId, positionId int64
	var err2 error
	
	if req.DepartmentId != "" {
		departmentId, err2 = util.StringToInt64(req.DepartmentId)
		if err2 != nil {
			return nil, err2
		}
	}
	if req.PositionId != "" {
		positionId, err2 = util.StringToInt64(req.PositionId)
		if err2 != nil {
			return nil, err2
		}
	}
	
	ret, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.SearchEmployeeDetail(l.ctx, &pb.SearchEmployeeDetailReq{
		Page:         req.Page,
		Limit:        req.Limit,
		Gender:       req.Gender,
		DepartmentId: departmentId,
		PositionId:   positionId,
		Salary:       float64(req.Salary),
		Name:         req.Name,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.SearchEmployeeResponse{
		Total: ret.Total,
	}
	for _, v := range ret.EmployeeNonSensitiveDetail {
		// 应该是关联查询，而非再调用
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

		resp.List = append(resp.List, &types.EmployeeDetail{
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
	return
}

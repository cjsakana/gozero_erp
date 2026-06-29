package position

import (
	"context"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"erp/app/hr/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPositionEmployeesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPositionEmployeesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPositionEmployeesLogic {
	return &GetPositionEmployeesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPositionEmployeesLogic) GetPositionEmployees(req *types.GetPositionEmployeesRequest) (resp *types.GetPositionEmployeesResponse, err error) {
	positionId, err := util.StringToInt64(req.PositionId)
	if err != nil {
		return nil, err
	}
	
	ret, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.SearchEmployeeDetail(l.ctx, &pb.SearchEmployeeDetailReq{
		Page:       req.Page,
		Limit:      req.Limit,
		PositionId: positionId,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.GetPositionEmployeesResponse{
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
		resp.List = append(resp.List, &types.EmployeeSimple{
			EmployeeId:     util.Int64ToString(v.Id),
			EmployeeNo:     v.EmployeeNo,
			EmployeeName:   v.Name,
			DepartmentId:   util.Int64ToString(departmentById.Department.Id),
			DepartmentName: departmentById.Department.Name,
		})
	}

	return
}

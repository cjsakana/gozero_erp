package employeeDetail

import (
	"context"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"erp/app/hr/rpc/pb"
	"erp/common/util"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmployeeDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetEmployeeDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmployeeDetailLogic {
	return &GetEmployeeDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEmployeeDetailLogic) GetEmployeeDetail(req *types.GetEmployeeDetailRequest) (resp *types.GetEmployeeDetailResponse, err error) {
	employeeId, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
		Id: employeeId,
	})
	if err != nil {
		return nil, err
	}

	departmentById, err2 := l.svcCtx.HrRPC.DepartmentZrpcClient.GetDepartmentById(l.ctx, &pb.GetDepartmentByIdReq{
		Id: ret.EmployeeNonSensitiveDetail.DepartmentId,
	})
	if err2 != nil {
		return nil, err2
	}

	positionById, err2 := l.svcCtx.HrRPC.PositionZrpcClient.GetPositionById(l.ctx, &pb.GetPositionByIdReq{
		Id: ret.EmployeeNonSensitiveDetail.PositionId,
	})
	if err2 != nil {
		return nil, err2
	}

	resp = &types.GetEmployeeDetailResponse{
		Employee: types.EmployeeDetail{
			Employee: types.Employee{
				Id:           util.Int64ToString(ret.EmployeeNonSensitiveDetail.Id),
				EmployeeNo:   ret.EmployeeNonSensitiveDetail.EmployeeNo,
				Name:         ret.EmployeeNonSensitiveDetail.Name,
				Gender:       ret.EmployeeNonSensitiveDetail.Gender,
				BirthDate:    ret.EmployeeNonSensitiveDetail.BirthDate,
				DepartmentId: util.Int64ToString(ret.EmployeeNonSensitiveDetail.DepartmentId),
				PositionId:   util.Int64ToString(ret.EmployeeNonSensitiveDetail.PositionId),
				Salary:       ret.EmployeeNonSensitiveDetail.Salary,
				HireDate:     ret.EmployeeNonSensitiveDetail.HireDate,
				LeaveDate:    ret.EmployeeNonSensitiveDetail.LeaveDate,
			},
			DepartmentName: departmentById.Department.Name,
			PositionName:   positionById.Position.Name,
		},
	}
	return
}

package employeeDetail

import (
	"context"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"erp/app/hr/rpc/pb"
	"erp/common/util"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateEmployeeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateEmployeeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateEmployeeLogic {
	return &UpdateEmployeeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateEmployeeLogic) UpdateEmployee(req *types.UpdateEmployeeRequest) (resp *types.EmptyResponse, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	departmentId, err := util.StringToInt64(req.DepartmentId)
	if err != nil {
		return nil, err
	}
	positionId, err := util.StringToInt64(req.PositionId)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.HrRPC.EmployeeDetailZrpcClient.UpdateEmployeeDetail(l.ctx, &pb.UpdateEmployeeDetailReq{
		Id:           id,
		Name:         req.Name,
		DepartmentId: departmentId,
		PositionId:   positionId,
		Salary:       req.Salary,
	})
	if err != nil {
		return nil, err
	}

	return
}

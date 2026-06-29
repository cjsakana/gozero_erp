package department

import (
	"context"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"erp/app/hr/rpc/pb"
	"erp/common/util"
	"erp/common/xtypes"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateDepartmentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDepartmentLogic {
	return &CreateDepartmentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateDepartmentLogic) CreateDepartment(req *types.CreateDepartmentRequest) (resp *types.CreateDepartmentResponse, err error) {
	employeeId, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}

	parentId, err := util.StringToInt64(req.ParentId)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.HrRPC.DepartmentZrpcClient.AddDepartment(l.ctx, &pb.AddDepartmentReq{
		Name:        req.Name,
		ParentId:    parentId,
		Code:        req.Code,
		ManagerId:   employeeId,
		ManagerNo:   req.ManagerNo,
		ManagerName: req.ManagerName,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.CreateDepartmentResponse{
		Id: util.Int64ToString(ret.Id),
	}
	return
}

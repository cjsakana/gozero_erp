package department

import (
	"context"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"erp/app/hr/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateDepartmentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDepartmentLogic {
	return &UpdateDepartmentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateDepartmentLogic) UpdateDepartment(req *types.UpdateDepartmentRequest) (resp *types.EmptyResponse, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	parentId, err := util.StringToInt64(req.ParentId)
	if err != nil {
		return nil, err
	}
	managerId, err := util.StringToInt64(req.ManagerId)
	if err != nil {
		return nil, err
	}
	
	_, err = l.svcCtx.HrRPC.DepartmentZrpcClient.UpdateDepartment(l.ctx, &pb.UpdateDepartmentReq{
		Id:          id,
		Name:        req.Name,
		ParentId:    parentId,
		Code:        req.Code,
		ManagerId:   managerId,
		ManagerNo:   req.ManagerNo,
		ManagerName: req.ManagerName,
	})
	if err != nil {
		return nil, err
	}

	return
}

package department

import (
	"context"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"erp/app/hr/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteDepartmentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDepartmentLogic {
	return &DeleteDepartmentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteDepartmentLogic) DeleteDepartment(req *types.DeleteDepartmentRequest) (resp *types.EmptyResponse, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	
	_, err = l.svcCtx.HrRPC.DepartmentZrpcClient.DelDepartment(l.ctx, &pb.DelDepartmentReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return
}

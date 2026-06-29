package resigned

import (
	"context"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"erp/app/hr/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateResignedApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateResignedApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateResignedApplicationLogic {
	return &UpdateResignedApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateResignedApplicationLogic) UpdateResignedApplication(req *types.UpdateResignedApplicationRequest) (resp *types.EmptyResponse, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.HrRPC.ResignedApplicationZrpcClient.UpdateResignedApplication(l.ctx, &pb.UpdateResignedApplicationReq{
		Id:        id,
		Reason:    req.Reason,
		LeaveDate: req.LeaveDate,
		Evidence:  req.Evidence,
	})
	if err != nil {
		return nil, err
	}
	return
}

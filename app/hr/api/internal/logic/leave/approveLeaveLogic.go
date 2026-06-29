package leave

import (
	"context"
	"erp/app/hr/rpc/pb"
	"erp/common/util"
	"time"

	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApproveLeaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApproveLeaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveLeaveLogic {
	return &ApproveLeaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApproveLeaveLogic) ApproveLeave(req *types.ApproveLeaveRequest) (resp *types.EmptyResponse, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.HrRPC.LeaveApplicationZrpcClient.UpdateLeaveApplication(l.ctx, &pb.UpdateLeaveApplicationReq{
		Id:            id,
		Status:        req.Status,
		ApproveTime:   time.Now().Unix(),
		ApproveRemark: req.ApproveRemark,
	})
	if err != nil {
		return nil, err
	}

	return
}

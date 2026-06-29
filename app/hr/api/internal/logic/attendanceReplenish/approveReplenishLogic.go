package attendanceReplenish

import (
	"context"
	"erp/app/hr/rpc/pb"
	"erp/common/util"
	"time"

	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApproveReplenishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApproveReplenishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveReplenishLogic {
	return &ApproveReplenishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApproveReplenishLogic) ApproveReplenish(req *types.ApproveReplenishRequest) (resp *types.EmptyResponse, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.HrRPC.AttendanceReplenishZrpcClient.UpdateAttendanceReplenish(l.ctx, &pb.UpdateAttendanceReplenishReq{
		Id:            id,
		Status:        req.Status,
		ApproveRemark: req.ApproveRemark,
		ApproveTime:   time.Now().Unix(),
	})
	if err != nil {
		return nil, err
	}
	return
}

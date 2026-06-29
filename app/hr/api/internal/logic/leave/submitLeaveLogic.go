package leave

import (
	"context"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"erp/app/hr/rpc/pb"
	"erp/common/util"
	"erp/common/xtypes"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubmitLeaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubmitLeaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitLeaveLogic {
	return &SubmitLeaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubmitLeaveLogic) SubmitLeave(req *types.SubmitLeaveRequest) (resp *types.SubmitLeaveResponse, err error) {
	employeeId, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}

	approverId, err := util.StringToInt64(req.ApproverId)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.HrRPC.LeaveApplicationZrpcClient.AddLeaveApplication(l.ctx, &pb.AddLeaveApplicationReq{
		EmployeeId: employeeId,
		Type:       req.Type,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		Reason:     req.Reason,
		Evidence:   req.Evidence,
		ApproverId: approverId,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.SubmitLeaveResponse{
		Id: util.Int64ToString(ret.Id),
	}
	return
}

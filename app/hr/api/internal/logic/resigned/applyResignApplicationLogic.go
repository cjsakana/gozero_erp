package resigned

import (
	"context"
	"erp/app/hr/rpc/pb"
	"erp/common/util"
	"erp/common/xtypes"

	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyResignApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApplyResignApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyResignApplicationLogic {
	return &ApplyResignApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApplyResignApplicationLogic) ApplyResignApplication(req *types.ApplyResignApplicationRequest) (resp *types.ApplyResignApplicationResponse, err error) {
	employeeId, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}

	approverId, err := util.StringToInt64(req.ApproverId)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.HrRPC.ResignedApplicationZrpcClient.AddResignedApplication(l.ctx, &pb.AddResignedApplicationReq{
		EmployeeId: employeeId,
		Reason:     req.Reason,
		LeaveDate:  req.LeaveDate,
		Evidence:   req.Evidence,
		ApproverId: approverId,
	})
	if err != nil {
		return nil, err
	}

	return &types.ApplyResignApplicationResponse{
		Id: util.Int64ToString(ret.Id),
	}, nil
}

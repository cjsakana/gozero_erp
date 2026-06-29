package attendanceReplenish

import (
	"context"
	"erp/app/hr/api/internal/code"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"erp/app/hr/rpc/pb"
	"erp/common/util"
	"erp/common/xtypes"
	"github.com/zeromicro/go-zero/core/logx"
)

type SubmitReplenishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubmitReplenishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitReplenishLogic {
	return &SubmitReplenishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubmitReplenishLogic) SubmitReplenish(req *types.SubmitReplenishRequest) (resp *types.SubmitReplenishResponse, err error) {
	employeeId, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}
	approverId, err := util.StringToInt64(req.ApproverId)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.HrRPC.AttendanceReplenishZrpcClient.AddAttendanceReplenish(l.ctx, &pb.AddAttendanceReplenishReq{
		EmployeeId:    employeeId,
		OriginalDate:  req.OriginalDate,
		ReplenishType: req.ReplenishType,
		ReplenishTime: req.ReplenishTime,
		Reason:        req.Reason,
		Evidence:      req.Evidence,
		ApproverId:    approverId,
	})
	if err != nil {
		return nil, code.SubmitReplenishFail
	}
	resp = &types.SubmitReplenishResponse{
		Id: util.Int64ToString(ret.Id),
	}
	return
}

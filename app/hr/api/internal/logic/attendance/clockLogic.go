package attendance

import (
	"context"
	"erp/app/hr/rpc/pb"
	"erp/common/util"
	"erp/common/xtypes"
	"time"

	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClockLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClockLogic {
	return &ClockLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClockLogic) Clock(req *types.ClockRequest) (resp *types.EmptyResponse, err error) {
	employeeId, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}

	// 通过工号查询员工ID
	employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
		Id: employeeId,
	})
	if err != nil {
		return nil, err
	}

	// 调用打卡RPC
	_, err = l.svcCtx.HrRPC.AttendanceRecordZrpcClient.Clock(l.ctx, &pb.ClockReq{
		EmployeeId: employeeDetail.EmployeeNonSensitiveDetail.Id,
		ClockTime:  time.Now().Unix(),
	})
	if err != nil {
		return nil, err
	}
	return &types.EmptyResponse{}, nil
}

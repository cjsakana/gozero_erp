package attendanceReplenish

import (
	"context"
	"erp/app/hr/rpc/pb"
	"erp/common/util"

	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetReplenishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetReplenishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetReplenishLogic {
	return &GetReplenishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetReplenishLogic) GetReplenish(req *types.GetReplenishRequest) (resp *types.GetReplenishResponse, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.HrRPC.AttendanceReplenishZrpcClient.GetAttendanceReplenishById(l.ctx, &pb.GetAttendanceReplenishByIdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	// 查询员工信息
	employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
		Id: ret.AttendanceReplenish.EmployeeId,
	})
	if err != nil {
		return nil, err
	}

	// 查询审批人信息（如果有）
	var approverNo string
	var approverName string
	if ret.AttendanceReplenish.ApproverId > 0 {
		approverDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
			Id: ret.AttendanceReplenish.ApproverId,
		})
		if err != nil {
			logx.Errorf("查询审批人信息失败: approverId=%d, err=%v", ret.AttendanceReplenish.ApproverId, err)
		} else {
			approverNo = approverDetail.EmployeeNonSensitiveDetail.EmployeeNo
			approverName = approverDetail.EmployeeNonSensitiveDetail.Name
		}
	}

	resp = &types.GetReplenishResponse{
		AttendanceReplenish: types.AttendanceReplenish{
			Id:            util.Int64ToString(ret.AttendanceReplenish.Id),
			EmployeeId:    util.Int64ToString(ret.AttendanceReplenish.EmployeeId),
			EmployeeNo:    employeeDetail.EmployeeNonSensitiveDetail.EmployeeNo,
			EmployeeName:  employeeDetail.EmployeeNonSensitiveDetail.Name,
			OriginalDate:  ret.AttendanceReplenish.OriginalDate,
			ApplyTime:     ret.AttendanceReplenish.ApplyTime,
			ReplenishType: ret.AttendanceReplenish.ReplenishType,
			ReplenishTime: ret.AttendanceReplenish.ReplenishTime,
			Reason:        ret.AttendanceReplenish.Reason,
			Evidence:      ret.AttendanceReplenish.Evidence,
			Status:        ret.AttendanceReplenish.Status,
			ApproverId:    util.Int64ToString(ret.AttendanceReplenish.ApproverId),
			ApproverNo:    approverNo,
			ApproverName:  approverName,
			ApproveTime:   ret.AttendanceReplenish.ApproveTime,
			ApproveRemark: ret.AttendanceReplenish.ApproveRemark,
		},
	}
	return
}

package leave

import (
	"context"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"erp/app/hr/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLeaveDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLeaveDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLeaveDetailLogic {
	return &GetLeaveDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLeaveDetailLogic) GetLeaveDetail(req *types.GetLeaveDetailRequest) (resp *types.GetLeaveDetailResponse, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.HrRPC.LeaveApplicationZrpcClient.GetLeaveApplicationById(l.ctx, &pb.GetLeaveApplicationByIdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	// 查询员工信息
	employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
		Id: ret.LeaveApplication.EmployeeId,
	})
	if err != nil {
		return nil, err
	}

	// 查询审批人信息（如果有）
	var approverNo string
	var approverName string
	if ret.LeaveApplication.ApproverId > 0 {
		approverDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
			Id: ret.LeaveApplication.ApproverId,
		})
		if err != nil {
			logx.Errorf("查询审批人信息失败: approverId=%d, err=%v", ret.LeaveApplication.ApproverId, err)
		} else {
			approverNo = approverDetail.EmployeeNonSensitiveDetail.EmployeeNo
			approverName = approverDetail.EmployeeNonSensitiveDetail.Name
		}
	}

	resp = &types.GetLeaveDetailResponse{
		Leave: types.LeaveApplication{
			Id:            util.Int64ToString(ret.LeaveApplication.Id),
			EmployeeId:    util.Int64ToString(ret.LeaveApplication.EmployeeId),
			EmployeeNo:    employeeDetail.EmployeeNonSensitiveDetail.EmployeeNo,
			EmployeeName:  employeeDetail.EmployeeNonSensitiveDetail.Name,
			Type:          ret.LeaveApplication.Type,
			StartTime:     ret.LeaveApplication.StartTime,
			EndTime:       ret.LeaveApplication.EndTime,
			Duration:      ret.LeaveApplication.Duration,
			Reason:        ret.LeaveApplication.Reason,
			Evidence:      ret.LeaveApplication.Evidence,
			Status:        ret.LeaveApplication.Status,
			ApproverId:    util.Int64ToString(ret.LeaveApplication.ApproverId),
			ApproverNo:    approverNo,
			ApproverName:  approverName,
			ApproveTime:   ret.LeaveApplication.ApproveTime,
			ApproveRemark: ret.LeaveApplication.ApproveRemark,
			CreatedAt:     ret.LeaveApplication.CreatedAt,
		},
	}

	return
}

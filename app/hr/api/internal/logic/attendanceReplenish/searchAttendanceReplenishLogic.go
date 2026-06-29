package attendanceReplenish

import (
	"context"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"erp/app/hr/rpc/client/employeedetail"
	"erp/app/hr/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchAttendanceReplenishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchAttendanceReplenishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchAttendanceReplenishLogic {
	return &SearchAttendanceReplenishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchAttendanceReplenishLogic) SearchAttendanceReplenish(req *types.SearchAttendanceReplenishRequest) (resp *types.SearchAttendanceReplenishResponse, err error) {
	employeeId, err := util.StringToInt64(req.EmployeeId)
	if err != nil {
		return nil, err
	}
	approverId, err := util.StringToInt64(req.ApproverId)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.HrRPC.AttendanceReplenishZrpcClient.SearchAttendanceReplenish(l.ctx, &pb.SearchAttendanceReplenishReq{
		Page:          req.Page,
		Limit:         req.Limit,
		EmployeeId:    employeeId,
		ReplenishType: req.ReplenishType,
		Reason:        req.Reason,
		Status:        req.Status,
		ApproverId:    approverId,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.SearchAttendanceReplenishResponse{
		Total: ret.Total,
	}

	// 批量查询员工和审批人信息
	employeeMap := make(map[int64]*employeedetail.EmployeeNonSensitiveDetail)
	for _, replenish := range ret.AttendanceReplenish {
		// 查询员工信息
		if _, ok := employeeMap[replenish.EmployeeId]; !ok {
			employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
				Id: replenish.EmployeeId,
			})
			if err != nil {
				logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", replenish.EmployeeId, err)
				continue
			}
			employeeMap[replenish.EmployeeId] = employeeDetail.EmployeeNonSensitiveDetail
		}

		// 查询审批人信息（如果有）
		if replenish.ApproverId > 0 {
			if _, ok := employeeMap[replenish.ApproverId]; !ok {
				approverDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
					Id: replenish.ApproverId,
				})
				if err != nil {
					logx.Errorf("查询审批人信息失败: approverId=%d, err=%v", replenish.ApproverId, err)
				} else {
					employeeMap[replenish.ApproverId] = approverDetail.EmployeeNonSensitiveDetail
				}
			}
		}

		resp.List = append(resp.List, &types.AttendanceReplenish{
			Id:            util.Int64ToString(replenish.Id),
			EmployeeId:    util.Int64ToString(replenish.EmployeeId),
			EmployeeNo:    employeeMap[replenish.EmployeeId].EmployeeNo,
			EmployeeName:  employeeMap[replenish.EmployeeId].Name,
			OriginalDate:  replenish.OriginalDate,
			ApplyTime:     replenish.ApplyTime,
			ReplenishType: replenish.ReplenishType,
			ReplenishTime: replenish.ReplenishTime,
			Reason:        replenish.Reason,
			Evidence:      replenish.Evidence,
			Status:        replenish.Status,
			ApproverId:    util.Int64ToString(replenish.ApproverId),
			ApproverNo:    employeeMap[replenish.ApproverId].EmployeeNo,
			ApproverName:  employeeMap[replenish.ApproverId].Name,
			ApproveTime:   replenish.ApproveTime,
			ApproveRemark: replenish.ApproveRemark,
		})
	}
	return
}

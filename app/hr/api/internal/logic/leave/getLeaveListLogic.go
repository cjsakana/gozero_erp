package leave

import (
	"context"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"erp/app/hr/rpc/client/employeedetail"
	"erp/app/hr/rpc/pb"
	"erp/common/util"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetLeaveListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLeaveListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLeaveListLogic {
	return &GetLeaveListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLeaveListLogic) GetLeaveList(req *types.SearchLeaveRequest) (resp *types.SearchLeaveResponse, err error) {
	employeeId, err := util.StringToInt64(req.EmployeeId)
	if err != nil {
		return nil, err
	}
	approverId, err := util.StringToInt64(req.ApproverId)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.HrRPC.LeaveApplicationZrpcClient.SearchLeaveApplication(l.ctx, &pb.SearchLeaveApplicationReq{
		Page:       req.Page,
		Limit:      req.Limit,
		EmployeeId: employeeId,
		Type:       req.Type,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		Reason:     req.Reason,
		Status:     req.Status,
		ApproverId: approverId,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.SearchLeaveResponse{
		Total: ret.Total,
	}

	// 批量查询员工和审批人信息
	employeeMap := make(map[int64]*employeedetail.EmployeeNonSensitiveDetail)
	for _, v := range ret.LeaveApplication {
		// 查询员工信息
		if _, ok := employeeMap[v.EmployeeId]; !ok {
			employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
				Id: v.EmployeeId,
			})
			if err != nil {
				logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", v.EmployeeId, err)
				continue
			}
			employeeMap[v.EmployeeId] = employeeDetail.EmployeeNonSensitiveDetail
		}

		// 查询审批人信息（如果有）
		if v.ApproverId > 0 {
			if _, ok := employeeMap[v.ApproverId]; !ok {
				approverDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
					Id: v.ApproverId,
				})
				if err != nil {
					logx.Errorf("查询审批人信息失败: approverId=%d, err=%v", v.ApproverId, err)
				} else {
					employeeMap[v.ApproverId] = approverDetail.EmployeeNonSensitiveDetail
				}
			}
		}

		resp.Application = append(resp.Application, &types.LeaveApplication{
			Id:            util.Int64ToString(v.Id),
			EmployeeId:    util.Int64ToString(v.EmployeeId),
			EmployeeNo:    employeeMap[v.EmployeeId].EmployeeNo,
			EmployeeName:  employeeMap[v.EmployeeId].Name,
			Type:          v.Type,
			StartTime:     v.StartTime,
			EndTime:       v.EndTime,
			Duration:      v.Duration,
			Reason:        v.Reason,
			Evidence:      v.Evidence,
			Status:        v.Status,
			ApproverId:    util.Int64ToString(v.ApproverId),
			ApproverNo:    employeeMap[v.ApproverId].EmployeeNo,
			ApproverName:  employeeMap[v.ApproverId].Name,
			ApproveTime:   v.ApproveTime,
			ApproveRemark: v.ApproveRemark,
			CreatedAt:     v.CreatedAt,
		})
	}

	return
}

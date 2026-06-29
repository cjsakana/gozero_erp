package leaveapplicationlogic

import (
	"context"
	"erp/app/hr/rpc/internal/types"
	"fmt"
	"time"

	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"

	"erp/app/hr/rpc/internal/code"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchLeaveApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchLeaveApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLeaveApplicationLogic {
	return &SearchLeaveApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchLeaveApplicationLogic) SearchLeaveApplication(in *pb.SearchLeaveApplicationReq) (*pb.SearchLeaveApplicationResp, error) {
	var startTime, endTime time.Time
	if in.StartTime > 0 {
		startTime = time.Unix(in.StartTime, 0)
	}
	if in.EndTime > 0 {
		endTime = time.Unix(in.EndTime, 0)
	}
	applications, total, err := l.svcCtx.LeaveApplicationModel.Search(l.ctx, &types.SearchLeaveApplicationParams{
		SearchCom: types.SearchCom{
			Page:  in.Page,
			Limit: in.Limit,
		},
		EmployeeId: in.EmployeeId, // 使用员工ID（新版主键）
		Type:       in.Type,
		StartTime:  startTime,
		EndTime:    endTime,
		Reason:     in.Reason,
		Status:     in.Status,
		ApproverId: in.ApproverId, // 使用审批人ID（新版主键）
	})
	if err != nil {
		fmt.Println("111111111111111111   ", err)
		return nil, code.GetLeaveFail
	}

	var leaveApplications []*pb.LeaveApplication
	for _, application := range applications {
		leaveApplications = append(leaveApplications, &pb.LeaveApplication{
			Id:            application.Id,
			EmployeeId:    application.EmployeeId, // 使用员工ID
			Type:          application.Type,
			StartTime:     application.StartTime.Unix(),
			EndTime:       application.EndTime.Unix(),
			Duration:      application.Duration,
			Reason:        application.Reason,
			Evidence:      application.Evidence.String,
			Status:        application.Status,
			ApproverId:    application.ApproverId.Int64, // 使用审批人ID
			ApproveTime:   application.ApproveTime.Time.Unix(),
			ApproveRemark: application.ApproveRemark.String,
			CreatedAt:     application.CreatedAt.Unix(),
		})
	}

	return &pb.SearchLeaveApplicationResp{
		Total:            total,
		LeaveApplication: leaveApplications,
	}, nil
}

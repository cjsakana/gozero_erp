package leaveapplicationlogic

import (
	"context"
	"database/sql"
	"erp/app/hr/rpc/internal/code"
	"erp/app/hr/rpc/internal/model"
	"time"

	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLeaveApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLeaveApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLeaveApplicationLogic {
	return &UpdateLeaveApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLeaveApplicationLogic) UpdateLeaveApplication(in *pb.UpdateLeaveApplicationReq) (*pb.UpdateLeaveApplicationResp, error) {
	startTime := time.Unix(in.StartTime, 0)
	endTime := time.Unix(in.EndTime, 0)
	duration := endTime.Sub(startTime)
	hours := duration.Hours()
	days := hours / 24

	err := l.svcCtx.LeaveApplicationModel.XUpdate(l.ctx, &model.LeaveApplication{
		Id:            in.Id,
		Type:          in.Type,
		StartTime:     startTime,
		EndTime:       endTime,
		Duration:      days,
		Reason:        in.Reason,
		Evidence:      sql.NullString{String: in.Evidence, Valid: in.Evidence != ""},
		Status:        in.Status,
		ApproverId:    sql.NullInt64{Int64: in.ApproverId, Valid: in.ApproverId != 0},
		ApproveTime:   sql.NullTime{Time: time.Unix(in.ApproveTime, 0), Valid: in.ApproveTime != 0},
		ApproveRemark: sql.NullString{String: in.ApproveRemark, Valid: in.ApproveRemark != ""},
	})
	if err != nil {
		return nil, code.ApproveLeaveFail
	}

	return &pb.UpdateLeaveApplicationResp{}, nil
}

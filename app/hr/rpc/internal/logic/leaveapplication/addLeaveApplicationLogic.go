package leaveapplicationlogic

import (
	"context"
	"database/sql"
	"erp/app/hr/rpc/internal/code"
	"erp/app/hr/rpc/internal/model"
	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"
	"erp/common/util"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type AddLeaveApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddLeaveApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLeaveApplicationLogic {
	return &AddLeaveApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------请假申请表-----------------------
func (l *AddLeaveApplicationLogic) AddLeaveApplication(in *pb.AddLeaveApplicationReq) (*pb.AddLeaveApplicationResp, error) {
	startTime := time.Unix(in.StartTime, 0)
	endTime := time.Unix(in.EndTime, 0)
	duration := endTime.Sub(startTime)
	hours := duration.Hours() // 小时数
	days := hours / 24        // 天数

	// 生成雪花ID
	id := util.GenerateSnowflake()
	_, err := l.svcCtx.LeaveApplicationModel.Insert(l.ctx, &model.LeaveApplication{
		Id:            id,
		EmployeeId:    in.EmployeeId,
		Type:          in.Type,
		StartTime:     startTime,
		EndTime:       endTime,
		Duration:      days,
		Reason:        in.Reason,
		Evidence:      sql.NullString{String: in.Evidence, Valid: in.Evidence != ""},
		Status:        1,
		ApproverId:    sql.NullInt64{Int64: in.ApproverId, Valid: in.ApproverId != 0},
		ApproveTime:   sql.NullTime{Valid: false},
		ApproveRemark: sql.NullString{Valid: false},
	})
	if err != nil {
		return nil, code.SubmitLeaveFail
	}
	return &pb.AddLeaveApplicationResp{
		Id: id,
	}, nil
}

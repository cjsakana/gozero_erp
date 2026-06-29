package attendancereplenishlogic

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

type UpdateAttendanceReplenishLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAttendanceReplenishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAttendanceReplenishLogic {
	return &UpdateAttendanceReplenishLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAttendanceReplenishLogic) UpdateAttendanceReplenish(in *pb.UpdateAttendanceReplenishReq) (*pb.UpdateAttendanceReplenishResp, error) {
	err := l.svcCtx.AttendanceReplenishModel.XUpdate(l.ctx, &model.AttendanceReplenish{
		Id:            in.Id,
		OriginalDate:  time.Unix(in.OriginalDate, 0),
		ReplenishType: in.ReplenishType,
		ReplenishTime: func() sql.NullTime {
			// 全天则不能使用
			if in.ReplenishType == 3 {
				return sql.NullTime{Valid: false}
			}
			return sql.NullTime{Time: time.Unix(in.ReplenishTime, 0), Valid: true}
		}(),
		Reason:        in.Reason,
		Evidence:      sql.NullString{String: in.Evidence, Valid: in.Evidence != ""},
		Status:        in.Status,
		ApproverId:    sql.NullInt64{Int64: in.ApproverId, Valid: in.ApproverId != 0},
		ApproveTime:   sql.NullTime{Time: time.Unix(in.ApproveTime, 0), Valid: in.ApproveTime != 0},
		ApproveRemark: sql.NullString{String: in.ApproveRemark, Valid: in.ApproveRemark != ""},
	})
	if err != nil {
		return nil, code.ApproveReplenishFail
	}

	return &pb.UpdateAttendanceReplenishResp{}, nil
}

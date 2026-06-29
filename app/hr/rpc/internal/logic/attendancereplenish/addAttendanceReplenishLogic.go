package attendancereplenishlogic

import (
	"context"
	"database/sql"
	"erp/app/hr/rpc/internal/code"
	"erp/app/hr/rpc/internal/model"
	"erp/common/util"
	"time"

	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddAttendanceReplenishLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddAttendanceReplenishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAttendanceReplenishLogic {
	return &AddAttendanceReplenishLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------补卡申请表-----------------------
func (l *AddAttendanceReplenishLogic) AddAttendanceReplenish(in *pb.AddAttendanceReplenishReq) (*pb.AddAttendanceReplenishResp, error) {
	// 生成雪花ID
	id := util.GenerateSnowflake()
	_, err := l.svcCtx.AttendanceReplenishModel.Insert(l.ctx, &model.AttendanceReplenish{
		Id:            id,
		EmployeeId:    in.EmployeeId,
		OriginalDate:  time.Unix(in.OriginalDate, 0),
		ApplyTime:     time.Now(),
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
		Status:        1,
		ApproverId:    sql.NullInt64{Int64: in.ApproverId, Valid: in.ApproverId != 0},
		ApproveTime:   sql.NullTime{Valid: false},
		ApproveRemark: sql.NullString{Valid: false},
	})
	if err != nil {
		return nil, code.SubmitReplenishFail
	}
	return &pb.AddAttendanceReplenishResp{
		Id: id,
	}, nil
}

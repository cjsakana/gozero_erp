package attendancereplenishlogic

import (
	"context"
	"erp/app/hr/rpc/internal/code"
	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

type GetAttendanceReplenishByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAttendanceReplenishByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAttendanceReplenishByIdLogic {
	return &GetAttendanceReplenishByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAttendanceReplenishByIdLogic) GetAttendanceReplenishById(in *pb.GetAttendanceReplenishByIdReq) (*pb.GetAttendanceReplenishByIdResp, error) {
	one, err := l.svcCtx.AttendanceReplenishModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, sqlc.ErrNotFound) {
			return nil, code.AttendanceNotFound
		}
		return nil, code.AttendanceNotFound
	}

	return &pb.GetAttendanceReplenishByIdResp{
		AttendanceReplenish: &pb.AttendanceReplenish{
			Id:            one.Id,
			EmployeeId:    one.EmployeeId,
			OriginalDate:  one.OriginalDate.Unix(),
			ApplyTime:     one.ApplyTime.Unix(),
			ReplenishType: one.ReplenishType,
			ReplenishTime: one.ReplenishTime.Time.Unix(),
			Reason:        one.Reason,
			Evidence:      one.Evidence.String,
			Status:        one.Status,
			ApproverId:    one.ApproverId.Int64,
			ApproveTime:   one.ApproveTime.Time.Unix(),
			ApproveRemark: one.ApproveRemark.String,
		},
	}, nil
}

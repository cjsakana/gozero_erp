package attendancerecordlogic

import (
	"context"
	"erp/app/hr/rpc/internal/code"
	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

type GetAttendanceRecordByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAttendanceRecordByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAttendanceRecordByIdLogic {
	return &GetAttendanceRecordByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAttendanceRecordByIdLogic) GetAttendanceRecordById(in *pb.GetAttendanceRecordByIdReq) (*pb.GetAttendanceRecordByIdResp, error) {
	one, err := l.svcCtx.AttendanceRecordModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, sqlc.ErrNotFound) {
			return nil, code.AttendanceNotFound
		}
		return nil, code.AttendanceNotFound
	}

	return &pb.GetAttendanceRecordByIdResp{
		AttendanceRecord: &pb.AttendanceRecord{
			Id:            one.Id,
			EmployeeId:    one.EmployeeId,
			Date:          one.Date.Unix(),
			ClockIn:       one.ClockIn.Time.Unix(),
			ClockOut:      one.ClockOut.Time.Unix(),
			IsAmMissing:   one.IsAmMissing == 1,
			IsLate:        one.IsLate == 1,
			IsPmMissing:   one.IsPmMissing == 1,
			IsEarlyLeave:  one.IsEarlyLeave == 1,
			WorkHours:     one.WorkHours,
			OvertimeHours: one.OvertimeHours,
			Remark:        one.Remark.String,
		},
	}, nil
}

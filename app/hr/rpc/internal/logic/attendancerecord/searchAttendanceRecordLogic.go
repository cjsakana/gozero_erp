package attendancerecordlogic

import (
	"context"
	types2 "erp/app/hr/rpc/internal/types"
	"time"

	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"

	"erp/app/hr/rpc/internal/code"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchAttendanceRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchAttendanceRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchAttendanceRecordLogic {
	return &SearchAttendanceRecordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchAttendanceRecordLogic) SearchAttendanceRecord(in *pb.SearchAttendanceRecordReq) (*pb.SearchAttendanceRecordResp, error) {
	var startDate, endDate time.Time
	if in.StartDate > 0 {
		startDate = time.Unix(in.StartDate, 0)
	}
	if in.EndDate > 0 {
		endDate = time.Unix(in.EndDate, 0)
	}
	// 如果只传了开始日期，默认按当天查询
	if !startDate.IsZero() && endDate.IsZero() {
		endDate = startDate
	}

	// 处理 optional 字段
	var isLate, isEarlyLeave, isAmMissing, isPmMissing bool
	if in.IsLate != nil {
		isLate = *in.IsLate
	}
	if in.IsEarlyLeave != nil {
		isEarlyLeave = *in.IsEarlyLeave
	}
	if in.IsAmMissing != nil {
		isAmMissing = *in.IsAmMissing
	}
	if in.IsPmMissing != nil {
		isPmMissing = *in.IsPmMissing
	}

	records, total, err := l.svcCtx.AttendanceRecordModel.Search(l.ctx, &types2.SearchAttendanceRecordParams{
		SearchCom: types2.SearchCom{
			Page:  in.Page,
			Limit: in.Limit,
		},
		EmployeeId:   in.EmployeeId,
		IsLate:       isLate,
		IsEarlyLeave: isEarlyLeave,
		IsAmMissing:  isAmMissing,
		IsPmMissing:  isPmMissing,
		Remark:       in.Remark,
		StartDate:    startDate,
		EndDate:      endDate,
	})
	if err != nil {
		return nil, code.SearchAttendanceFail
	}

	var pbRecords []*pb.AttendanceRecord
	for _, record := range records {
		pbRecords = append(pbRecords, &pb.AttendanceRecord{
			Id:            record.Id,
			EmployeeId:    record.EmployeeId, // 使用员工ID
			Date:          record.Date.Unix(),
			ClockIn:       record.ClockIn.Time.Unix(),
			ClockOut:      record.ClockOut.Time.Unix(),
			IsAmMissing:   record.IsAmMissing == 1,
			IsLate:        record.IsLate == 1,
			IsPmMissing:   record.IsPmMissing == 1,
			IsEarlyLeave:  record.IsEarlyLeave == 1,
			WorkHours:     record.WorkHours,
			OvertimeHours: record.OvertimeHours,
			Remark:        record.Remark.String,
		})
	}

	return &pb.SearchAttendanceRecordResp{
		Total:            total,
		AttendanceRecord: pbRecords,
	}, nil
}

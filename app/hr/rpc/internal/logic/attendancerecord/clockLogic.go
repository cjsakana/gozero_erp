package attendancerecordlogic

import (
	"context"
	"database/sql"
	"erp/app/hr/rpc/internal/model"
	"erp/common/util"
	"erp/common/xtime"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"time"

	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewClockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClockLogic {
	return &ClockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------考勤记录表-----------------------
func (l *ClockLogic) Clock(in *pb.ClockReq) (*pb.ClockResp, error) {
	// 初始化状态标志
	var isAmMissing, isLate, isPmMissing, isEarlyLeave int64 = 0, 0, 0, 0
	var remarks string = ""
	now := time.Now()
	loc, _ := time.LoadLocation("Asia/Shanghai")
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	noon12 := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, loc)

	clockType := 1 // 打卡类型：1-上午 2-下午
	if in.ClockTime > noon12.Unix() {
		clockType = 2
	}

	// 打卡（新建/更新）
	date, err := l.svcCtx.AttendanceRecordModel.FindOneByEmployeeIdDate(l.ctx, in.EmployeeId, today)
	switch err {
	case nil:
		// 修改：已经存在即上班打卡了，判断早退？
		if xtime.IsEarlyShanghai(in.GetClockTime()) {
			isEarlyLeave = 1
			remarks = fmt.Sprintf("早退 %d分钟",
				(xtime.GetStandardClockOutTime(in.ClockTime, loc)-in.ClockTime)/60)
		}
		clockOut := time.Unix(in.ClockTime, 0)
		workDuration := clockOut.Sub(date.ClockIn.Time)
		workHours := workDuration.Hours()

		// 计算超出8小时的加班时长
		var overtimeHours float64
		if workHours > 8.0 {
			overtimeHours = workHours - 8.0
		}

		err = l.svcCtx.AttendanceRecordModel.XUpdate(l.ctx, &model.AttendanceRecord{
			Id:            date.Id,
			EmployeeId:    date.EmployeeId,
			Date:          date.Date,
			ClockOut:      sql.NullTime{Time: clockOut, Valid: true},
			IsEarlyLeave:  isEarlyLeave,
			WorkHours:     workHours,
			OvertimeHours: overtimeHours,
			Remark:        sql.NullString{String: remarks, Valid: remarks != ""},
		})
		if err != nil {
			return nil, err
		}
	case sqlc.ErrNotFound:
		// 创建新记录
		if clockType == 1 {
			// 上午打卡：判断迟到？
			if xtime.IsLateShanghai(in.ClockTime) {
				isLate = 1
				remarks = fmt.Sprintf("迟到 %d分钟",
					(in.ClockTime-xtime.GetStandardClockInTime(in.ClockTime, loc))/60)
			}
		} else {
			// 下午打卡：上午缺卡
			isAmMissing = 1
			// 判断早退？
			if xtime.IsEarlyShanghai(in.GetClockTime()) {
				isEarlyLeave = 1
				if remarks != "" {
					remarks += "，"
				}
				remarks += fmt.Sprintf("早退 %d分钟",
					(xtime.GetStandardClockOutTime(in.ClockTime, loc)-in.ClockTime)/60)
			}
		}

		// 生成雪花ID
		id := util.GenerateSnowflake()
		record := &model.AttendanceRecord{
			Id:            id,
			EmployeeId:    in.EmployeeId, // 使用员工ID（新版主键）
			Date:          today,
			ClockIn:       sql.NullTime{Time: time.Unix(in.ClockTime, 0), Valid: clockType == 1},
			ClockOut:      sql.NullTime{Time: time.Unix(in.ClockTime, 0), Valid: clockType == 2},
			IsAmMissing:   isAmMissing,
			IsLate:        isLate,
			IsPmMissing:   isPmMissing,
			IsEarlyLeave:  isEarlyLeave,
			WorkHours:     0,
			OvertimeHours: 0,
			Remark:        sql.NullString{String: remarks, Valid: remarks != ""},
		}

		// 插入数据库
		_, err := l.svcCtx.AttendanceRecordModel.Insert(l.ctx, record)
		if err != nil {
			return nil, err
		}
	default:
		return nil, err
	}

	return &pb.ClockResp{}, nil
}

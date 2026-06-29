package attendancerecordlogic

import (
	"context"
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"

	"erp/app/hr/rpc/internal/code"
	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AttendanceCheckerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAttendanceCheckerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AttendanceCheckerLogic {
	return &AttendanceCheckerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AttendanceCheckerLogic) AttendanceChecker(in *pb.AttendanceCheckerReq) (*pb.AttendanceCheckerResp, error) {
	// 查询缺下班打卡的记录
	records, err := l.svcCtx.AttendanceRecordModel.FindMissingClockOut(l.ctx, time.Now())
	if err != nil {
		return nil, code.GetAttendanceFail
	}

	lock := redis.NewRedisLock(l.svcCtx.BizRedis, "attendanceCheckerLock")
	// 设置过期时间
	lock.SetExpire(30 * 60)
	// 尝试获取锁
	acquire, err := lock.Acquire()

	switch {
	case err != nil:
		// deal err
		logx.Errorf(err.Error())
	case acquire:
		// 获取到锁
		defer lock.Release() // 释放锁
		// 业务逻辑
		// 批量更新状态
		for _, record := range records {
			// 只更新未设置下班打卡的记录
			if !record.ClockOut.Valid {
				// 标记缺下午卡
				record.IsPmMissing = 1
				remarkText := "缺下班打卡"
				if record.Remark.Valid && record.Remark.String != "" {
					remarkText = record.Remark.String + "，" + remarkText
				}
				record.Remark = sql.NullString{
					String: remarkText,
					Valid:  true,
				}
				if err := l.svcCtx.AttendanceRecordModel.Update(l.ctx, record); err != nil {
					logx.WithContext(l.ctx).Errorf("更新记录失败 ID:%d %v", record.Id, err)
					continue
				}
			}
		}

		logx.Infof("已完成缺下班打卡检查，共处理%d条记录", len(records))

	case !acquire:
		// 没有拿到锁 wait?
	}

	return &pb.AttendanceCheckerResp{}, nil
}

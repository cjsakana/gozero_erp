package xtime

import (
	"strings"
)

const (
	LateMorning     = 1 << 0 // 00000001
	AbsentMorning   = 1 << 1 // 00000010
	EarlyAfternoon  = 1 << 2 // 00000100
	AbsentAfternoon = 1 << 3 // 00001000
)

func ExplainStatus(status int64) string {
	desc := []string{}

	if status&LateMorning != 0 {
		desc = append(desc, "上午迟到")
	}
	if status&AbsentMorning != 0 {
		desc = append(desc, "缺上午卡")
	}
	if status&EarlyAfternoon != 0 {
		desc = append(desc, "下午早退")
	}
	if status&AbsentAfternoon != 0 {
		desc = append(desc, "缺下午卡")
	}

	if len(desc) == 0 {
		return "全天正常"
	}
	return strings.Join(desc, "+")
}

//func CalculateAttendanceStatus(clockIn, clockOut time.Time) byte {
//	var status byte
//
//	// 基础时间配置
//	loc := time.FixedZone("CST", 8 * 3600) // 上海时区
//	workStart := time.Date(clockIn.Year(), clockIn.Month(), clockIn.Day(), 9, 0, 0, 0, loc)
//	workEnd := time.Date(clockOut.Year(), clockOut.Month(), clockOut.Day(), 18, 0, 0, 0, loc)
//
//	// 上午状态判断
//	if clockIn.IsZero() {
//		status |= AbsentMorning
//	} else if clockIn.After(workStart.Add(30*time.Minute)) {
//		status |= LateMorning
//	}
//
//	// 下午状态判断
//	if clockOut.IsZero() {
//		status |= AbsentAfternoon
//	} else if clockOut.Before(workEnd.Add(-60*time.Minute)) {
//		status |= EarlyAfternoon
//	}
//
//	return status
//}

//// 查询所有上午迟到的记录
//func FindLateMorningRecords(ctx context.Context) ([]*AttendanceRecord, error) {
//	return l.svcCtx.AttendanceRecordModel.FindWhere(ctx, "status & ? > 0", LateMorning)
//}
//
//// 查询缺上午卡但下午正常的记录
//func FindAbsentMorningOnly(ctx context.Context) ([]*AttendanceRecord, error) {
//	return l.svcCtx.AttendanceRecordModel.FindWhere(ctx,
//		"status & ? > 0 AND status & ? = 0",
//		AbsentMorning, AbsentAfternoon)
//}

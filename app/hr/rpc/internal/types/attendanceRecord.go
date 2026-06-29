package types

import "time"

type (
	// SearchAttendanceRecordParams 考勤记录搜索参数
	SearchAttendanceRecordParams struct {
		SearchCom
		EmployeeId   int64     // 员工ID
		IsLate       bool      // 是否迟到
		IsEarlyLeave bool      // 是否早退
		IsAmMissing  bool      // 是否缺上午卡
		IsPmMissing  bool      // 是否缺下午卡
		Remark       string    // 备注（模糊查询）
		StartDate    time.Time // 开始日期
		EndDate      time.Time // 结束日期
	}
)

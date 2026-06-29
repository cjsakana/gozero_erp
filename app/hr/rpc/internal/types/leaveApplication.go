package types

import "time"

type (
	// SearchLeaveApplicationParams 请假申请搜索参数
	SearchLeaveApplicationParams struct {
		SearchCom
		EmployeeId int64     // 员工ID（新版主键）
		Type       int64     // 请假类型
		StartTime  time.Time // 开始时间
		EndTime    time.Time // 结束时间
		Reason     string    // 请假原因
		Status     int64     // 状态
		ApproverId int64     // 审批人ID（新版主键）
	}
)

package types

import "time"

type (
	// SearchResignedApplicationParams 离职申请搜索参数
	SearchResignedApplicationParams struct {
		SearchCom
		EmployeeId     int64     // 申请人ID
		StartLeaveDate time.Time // 离职开始日期
		EndLeaveDate   time.Time // 离职结束日期
		Status         int64     // 状态
		ApproverId     int64     // 审批人ID（新版主键）
	}
)

package types

type (
	// SearchReplenishParams 补卡申请搜索参数
	SearchReplenishParams struct {
		SearchCom
		EmployeeId    int64  // 员工ID（新版主键）
		ReplenishType int64  // 补卡类型
		Reason        string // 补卡原因
		Status        int64  // 状态
		ApproverId    int64  // 审批人ID（新版主键）
	}
)

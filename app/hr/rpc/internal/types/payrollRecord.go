package types

import "time"

// SearchPayrollRecordParams 薪资记录搜索参数
type SearchPayrollRecordParams struct {
	SearchCom
	EmployeeId          int64     // 员工ID
	Status              int64     // 状态
	Description         string    // 描述
	CalculatedBy        int64     // 核算人ID
	StartCalculatedDate time.Time // 核算开始日期
	EndCalculatedDate   time.Time // 核算结束日期
	StartPaymentDate    time.Time // 发放开始日期
	EndPaymentDate      time.Time // 发放结束日期
	PaymentMonth        time.Time // 薪资月份
}

// BulkAddPayrollRecordErrItem 批量添加薪资记录错误项
type BulkAddPayrollRecordErrItem struct {
	EmployeeId int64 // 员工ID
	Success    bool  // 是否成功
	Err        error // 错误信息
}

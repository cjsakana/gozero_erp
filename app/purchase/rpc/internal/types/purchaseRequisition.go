package types

type RequisitionDetailParam struct {
	Id           int64   // 雪花ID
	ProductId    int64
	ProductName  string
	CategoryType int64
	Quantity     float64
	UnitPrice    float64
	Amount       float64
	Remark       string
}

type CreateRequisitionWithDetailsParam struct {
	RequisitionNo string
	DepartmentId  int64
	ApplicantId   int64
	ApproverId    int64
	RequestDate   int64
	TotalAmount   float64
	Status        int64
	Details       []RequisitionDetailParam
}

type SearchRequisitionParams struct {
	SearchComm
	RequisitionNo string
	DepartmentId  int64
	ApplicantId   int64
	ApproverId    int64
	Status        int64
}

type ApproveRequisitionParam struct {
	Id            int64
	ApproveTime   int64
	ApproveRemark string
	TargetStatus  int64
}

// 更新采购申请参数
type UpdateRequisitionParam struct {
	Id            int64
	DepartmentId  *int64   // 使用指针表示可选字段
	ApplicantId   *int64
	RequestDate   *int64
	TotalAmount   *float64
	Status        *int64
	ApproverId    *int64
	ApproveTime   *int64
	ApproveRemark *string
}

// 更新采购申请明细参数
type UpdateRequisitionDetailParam struct {
	Id           int64
	ProductId    *int64   // 使用指针表示可选字段
	ProductName  *string
	CategoryType *int64
	Quantity     *float64
	UnitPrice    *float64
	Amount       *float64
	Remark       *string
}

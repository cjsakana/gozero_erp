package types

type OrderDetailParam struct {
	Id           int64   // 雪花ID
	ProductId    int64
	ProductName  string
	CategoryType int64
	Quantity     float64
	UnitPrice    float64
	Amount       float64
	Remark       string
}

type CreateOrderWithDetailsParam struct {
	OrderNo      string
	SupplierId   int64
	OrderDate    int64
	ExpectedDate int64
	TotalAmount  float64
	Status       int64
	PurchaserId  int64
	Details      []OrderDetailParam
}

type CreateOrderFromRequisitionParam struct {
	RequisitionId int64
	OrderNo       string
	SupplierId    int64
	OrderDate     int64
	ExpectedDate  int64
	PurchaserId   int64
	Details       []OrderDetailParam // 可选覆盖
}

type SearchOrderParams struct {
	SearchComm
	OrderNo    string
	SupplierId int64
	Status     int64
}

// 更新采购订单参数
type UpdateOrderParam struct {
	Id           int64
	SupplierId   *int64   // 使用指针表示可选字段
	OrderDate    *int64
	ExpectedDate *int64
	TotalAmount  *float64
	Status       *int64
	PurchaserId  *int64
	ContractUrl  *string
}

// 更新采购订单明细参数
type UpdateOrderDetailParam struct {
	Id           int64
	ProductId    *int64   // 使用指针表示可选字段
	ProductName  *string
	CategoryType *int64
	Quantity     *float64
	UnitPrice    *float64
	Amount       *float64
	ReceivedQty  *float64
	Remark       *string
}

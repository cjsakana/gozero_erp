package types

type ReceiptDetailParam struct {
	Id           int64
	ProductId    int64
	ProductName  string
	CategoryType int64
	Quantity     float64
	UnitPrice    float64
	Amount       float64
	BatchId      int64
}

type CreateReceiptWithDetailsParam struct {
	ReceiptNo     string
	OrderId       int64
	WarehouseId   int64
	ReceiptDate   int64
	TotalQuantity float64
	TotalAmount   float64
	Status        int64
	CreatedBy     int64
	Details       []ReceiptDetailParam
}

type CreateReceiptFromOrderParam struct {
	OrderId     int64
	ReceiptNo   string
	WarehouseId int64
	ReceiptDate int64
	CreatedBy   int64
	Details     []ReceiptDetailParam // 可选覆盖
}

type SearchReceiptParams struct {
	SearchComm
	ReceiptNo   string
	OrderId     int64
	WarehouseId int64
}

// 更新采购入库单参数
type UpdateReceiptParam struct {
	Id            int64
	OrderId       *int64   // 使用指针表示可选字段
	WarehouseId   *int64
	ReceiptDate   *int64
	TotalQuantity *float64
	TotalAmount   *float64
	Status        *int64
	CreatedBy     *int64
}

// 更新采购入库明细参数
type UpdateReceiptDetailParam struct {
	Id           int64
	ProductId    *int64   // 使用指针表示可选字段
	ProductName  *string
	CategoryType *int64
	Quantity     *float64
	UnitPrice    *float64
	Amount       *float64
	BatchId      *int64
}

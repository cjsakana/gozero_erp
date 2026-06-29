package types

import "time"

type DeliveryDetailParam struct {
	Id          int64
	DeliveryId  int64
	ProductId   int64
	ProductName string
	Unit        string
	Quantity    float64
	UnitPrice   float64
	Amount      float64
	BatchId     int64
}

type AddSalesDeliveryParam struct {
	Id            int64
	DeliveryNo    string
	OrderId       int64
	WarehouseId   int64
	DeliveryDate  int64
	TotalQuantity float64
	TotalAmount   float64
	CreatedBy     int64
	Details       []DeliveryDetailParam
}

type OutboundDetailParam struct {
	Id        int64 // DeliveryId
	Quantity  float64
	ProductId int64
	BatchId   int64
}

type OutboundParam struct {
	Id        int64
	CreatedBy int64
	Items     []OutboundDetailParam
}

type SearchDeliveryParams struct {
	SearchComm
	DeliveryNo   string
	OrderId      int64
	WarehouseId  int64
	DeliveryDate time.Time
	Status       int64
}

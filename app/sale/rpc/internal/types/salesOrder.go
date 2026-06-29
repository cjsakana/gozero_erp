package types

import "time"

type SalesOrderDetailParam struct {
	Id           int64
	OrderId      int64
	ProductId    int64
	ProductName  string
	Unit         string
	Quantity     float64
	UnitPrice    float64
	Amount       float64
	DeliveredQty float64
	Remark       string
}

type AddSalesOrderParam struct {
	Id           int64
	OrderNo      string
	CustomerId   int64
	OrderDate    int64
	PromisedDate int64
	TotalAmount  float64
	Status       int64
	SalesmanId   int64
	ContractUrl  string
	Details      []SalesOrderDetailParam
}

type SearchOrderParams struct {
	SearchComm
	OrderNo        string
	CustomerId     int64
	Status         int64
	SalesmanId     int64
	StartOrderDate time.Time
	EndOrderDate   time.Time
}

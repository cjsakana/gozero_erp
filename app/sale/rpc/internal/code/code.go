package code

import "erp/common/xcode"

var (
	// 销售订单相关错误 108101-108110
	SalesOrderNotFound      = xcode.New(108101, "销售订单不存在")
	SalesOrderAlreadyExists = xcode.New(108102, "销售订单已存在")
	CreateSalesOrderFail    = xcode.New(108103, "创建销售订单失败")
	UpdateSalesOrderFail    = xcode.New(108104, "更新销售订单失败")
	UpdateOrderStatusFail   = xcode.New(108105, "更新销售订单状态失败")
	OrderStatusInvalid      = xcode.New(108106, "销售订单状态无效")
	GetSalesOrderFail       = xcode.New(108107, "获取销售订单失败")
	OrderAlreadyCompleted   = xcode.New(108108, "销售订单已完成")

	// 销售发货相关错误 108111-108120
	SalesDeliveryNotFound        = xcode.New(108111, "销售发货单不存在")
	SalesDeliveryAlreadyExists   = xcode.New(108112, "销售发货单已存在")
	CreateDeliveryFail           = xcode.New(108113, "创建销售发货单失败")
	UpdateDeliveryFail           = xcode.New(108114, "更新销售发货单失败")
	OutboundFail                 = xcode.New(108115, "出库失败")
	GetDeliveryFail              = xcode.New(108116, "获取销售发货单失败")
	DeliveryStatusInvalid        = xcode.New(108117, "销售发货单状态无效")
	StockInsufficientForOutbound = xcode.New(108118, "库存不足，无法出库")
)

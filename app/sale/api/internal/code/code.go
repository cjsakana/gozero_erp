package code

import "erp/common/xcode"

var (
	// 文件上传相关错误 108001-108010
	PutBucketErr = xcode.New(108001, "上传bucket失败")

	// 销售订单相关错误 108011-108020
	SalesOrderNotFound      = xcode.New(108011, "销售订单不存在")
	SalesOrderAlreadyExists = xcode.New(108012, "销售订单已存在")
	CreateSalesOrderFail    = xcode.New(108013, "创建销售订单失败")
	UpdateSalesOrderFail    = xcode.New(108014, "更新销售订单失败")
	UpdateOrderStatusFail   = xcode.New(108015, "更新销售订单状态失败")
	OrderStatusInvalid      = xcode.New(108016, "销售订单状态无效")
	GetSalesOrderFail       = xcode.New(108017, "获取销售订单失败")
	OrderAlreadyCompleted   = xcode.New(108018, "销售订单已完成")
	UploadContractFail      = xcode.New(108019, "上传合同文件失败")

	// 销售发货相关错误 108021-108030
	SalesDeliveryNotFound        = xcode.New(108021, "销售发货单不存在")
	SalesDeliveryAlreadyExists   = xcode.New(108022, "销售发货单已存在")
	CreateDeliveryFail           = xcode.New(108023, "创建销售发货单失败")
	UpdateDeliveryFail           = xcode.New(108024, "更新销售发货单失败")
	OutboundFail                 = xcode.New(108025, "出库失败")
	GetDeliveryFail              = xcode.New(108026, "获取销售发货单失败")
	DeliveryStatusInvalid        = xcode.New(108027, "销售发货单状态无效")
	StockInsufficientForOutbound = xcode.New(108028, "库存不足，无法出库")
)

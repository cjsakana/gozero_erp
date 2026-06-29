package code

import "erp/common/xcode"

var (
	// 采购申请相关错误 107101-107110
	PurchaseRequisitionNotFound      = xcode.New(107101, "采购申请不存在")
	PurchaseRequisitionAlreadyExists = xcode.New(107102, "采购申请已存在")
	CreateRequisitionFail            = xcode.New(107103, "创建采购申请失败")
	UpdateRequisitionFail            = xcode.New(107104, "更新采购申请失败")
	ApproveRequisitionFail           = xcode.New(107105, "审批采购申请失败")
	RequisitionStatusInvalid         = xcode.New(107106, "采购申请状态无效")
	GetRequisitionFail               = xcode.New(107107, "获取采购申请失败")
	RequisitionAlreadyApproved       = xcode.New(107108, "采购申请已审批")

	// 采购订单相关错误 107111-107120
	PurchaseOrderNotFound      = xcode.New(107111, "采购订单不存在")
	PurchaseOrderAlreadyExists = xcode.New(107112, "采购订单已存在")
	CreateOrderFail            = xcode.New(107113, "创建采购订单失败")
	UpdateOrderFail            = xcode.New(107114, "更新采购订单失败")
	CancelOrderFail            = xcode.New(107115, "取消采购订单失败")
	OrderStatusInvalid         = xcode.New(107116, "采购订单状态无效")
	GetOrderFail               = xcode.New(107117, "获取采购订单失败")
	OrderAlreadyCompleted      = xcode.New(107118, "采购订单已完成")

	// 采购收货相关错误 107121-107130
	PurchaseReceiptNotFound      = xcode.New(107121, "采购收货单不存在")
	PurchaseReceiptAlreadyExists = xcode.New(107122, "采购收货单已存在")
	CreateReceiptFail            = xcode.New(107123, "创建采购收货单失败")
	UpdateReceiptFail            = xcode.New(107124, "更新采购收货单失败")
	GetReceiptFail               = xcode.New(107125, "获取采购收货单失败")
	ReceiptStatusInvalid         = xcode.New(107126, "采购收货单状态无效")
)

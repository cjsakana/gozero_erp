package code

import "erp/common/xcode"

var (
	PutBucketErr = xcode.New(108001, "上传bucket失败")
	// 采购申请相关错误 107001-107010
	PurchaseRequisitionNotFound      = xcode.New(107001, "采购申请不存在")
	PurchaseRequisitionAlreadyExists = xcode.New(107002, "采购申请已存在")
	CreateRequisitionFail            = xcode.New(107003, "创建采购申请失败")
	UpdateRequisitionFail            = xcode.New(107004, "更新采购申请失败")
	ApproveRequisitionFail           = xcode.New(107005, "审批采购申请失败")
	RequisitionStatusInvalid         = xcode.New(107006, "采购申请状态无效")
	GetRequisitionFail               = xcode.New(107007, "获取采购申请失败")
	RequisitionAlreadyApproved       = xcode.New(107008, "采购申请已审批")

	// 采购订单相关错误 107011-107020
	PurchaseOrderNotFound      = xcode.New(107011, "采购订单不存在")
	PurchaseOrderAlreadyExists = xcode.New(107012, "采购订单已存在")
	CreateOrderFail            = xcode.New(107013, "创建采购订单失败")
	UpdateOrderFail            = xcode.New(107014, "更新采购订单失败")
	CancelOrderFail            = xcode.New(107015, "取消采购订单失败")
	OrderStatusInvalid         = xcode.New(107016, "采购订单状态无效")
	GetOrderFail               = xcode.New(107017, "获取采购订单失败")
	OrderAlreadyCompleted      = xcode.New(107018, "采购订单已完成")

	// 采购收货相关错误 107021-107030
	PurchaseReceiptNotFound      = xcode.New(107021, "采购收货单不存在")
	PurchaseReceiptAlreadyExists = xcode.New(107022, "采购收货单已存在")
	CreateReceiptFail            = xcode.New(107023, "创建采购收货单失败")
	UpdateReceiptFail            = xcode.New(107024, "更新采购收货单失败")
	GetReceiptFail               = xcode.New(107025, "获取采购收货单失败")
	ReceiptStatusInvalid         = xcode.New(107026, "采购收货单状态无效")
)

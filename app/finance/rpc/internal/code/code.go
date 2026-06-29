package code

import "erp/common/xcode"

var (
	// 固定资产相关错误 109101-109110
	FixedAssetNotFound      = xcode.New(109101, "固定资产不存在")
	FixedAssetAlreadyExists = xcode.New(109102, "固定资产已存在")
	AssetNoDuplicate        = xcode.New(109103, "资产编号重复")
	AddFixedAssetFail       = xcode.New(109104, "添加固定资产失败")
	UpdateFixedAssetFail    = xcode.New(109105, "更新固定资产失败")
	DeleteFixedAssetFail    = xcode.New(109106, "删除固定资产失败")
	GetFixedAssetFail       = xcode.New(109107, "获取固定资产失败")
	AssetInUse              = xcode.New(109108, "固定资产正在使用中，无法删除")

	// 付款记录相关错误 109111-109120
	PaymentRecordNotFound      = xcode.New(109111, "付款记录不存在")
	PaymentRecordAlreadyExists = xcode.New(109112, "付款记录已存在")
	AddPaymentRecordFail       = xcode.New(109113, "添加付款记录失败")
	UpdatePaymentRecordFail    = xcode.New(109114, "更新付款记录失败")
	DeletePaymentRecordFail    = xcode.New(109115, "删除付款记录失败")
	GetPaymentRecordFail       = xcode.New(109116, "获取付款记录失败")
	PaymentAmountInvalid       = xcode.New(109117, "付款金额无效")

	// 收款记录相关错误 109121-109130
	ReceiptRecordNotFound      = xcode.New(109121, "收款记录不存在")
	ReceiptRecordAlreadyExists = xcode.New(109122, "收款记录已存在")
	AddReceiptRecordFail       = xcode.New(109123, "添加收款记录失败")
	UpdateReceiptRecordFail    = xcode.New(109124, "更新收款记录失败")
	DeleteReceiptRecordFail    = xcode.New(109125, "删除收款记录失败")
	GetReceiptRecordFail       = xcode.New(109126, "获取收款记录失败")
	ReceiptAmountInvalid       = xcode.New(109127, "收款金额无效")

	// 工资发放相关错误 109131-109140
	SalaryPaymentNotFound      = xcode.New(109131, "工资发放记录不存在")
	SalaryPaymentAlreadyExists = xcode.New(109132, "工资发放记录已存在")
	AddSalaryPaymentFail       = xcode.New(109133, "添加工资发放记录失败")
	UpdateSalaryPaymentFail    = xcode.New(109134, "更新工资发放记录失败")
	DeleteSalaryPaymentFail    = xcode.New(109135, "删除工资发放记录失败")
	GetSalaryPaymentFail       = xcode.New(109136, "获取工资发放记录失败")
	SalaryAmountInvalid        = xcode.New(109137, "工资金额无效")
	PaymentStatusInvalid       = xcode.New(109138, "发放状态无效")
)

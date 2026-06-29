package code

import "erp/common/xcode"

var (
	// 固定资产相关错误 109001-109010
	FixedAssetNotFound      = xcode.New(109001, "固定资产不存在")
	FixedAssetAlreadyExists = xcode.New(109002, "固定资产已存在")
	AssetNoDuplicate        = xcode.New(109003, "资产编号重复")
	AddFixedAssetFail       = xcode.New(109004, "添加固定资产失败")
	UpdateFixedAssetFail    = xcode.New(109005, "更新固定资产失败")
	DeleteFixedAssetFail    = xcode.New(109006, "删除固定资产失败")
	GetFixedAssetFail       = xcode.New(109007, "获取固定资产失败")
	AssetInUse              = xcode.New(109008, "固定资产正在使用中，无法删除")

	// 付款记录相关错误 109011-109020
	PaymentRecordNotFound      = xcode.New(109011, "付款记录不存在")
	PaymentRecordAlreadyExists = xcode.New(109012, "付款记录已存在")
	AddPaymentRecordFail       = xcode.New(109013, "添加付款记录失败")
	UpdatePaymentRecordFail    = xcode.New(109014, "更新付款记录失败")
	DeletePaymentRecordFail    = xcode.New(109015, "删除付款记录失败")
	GetPaymentRecordFail       = xcode.New(109016, "获取付款记录失败")
	PaymentAmountInvalid       = xcode.New(109017, "付款金额无效")

	// 收款记录相关错误 109021-109030
	ReceiptRecordNotFound      = xcode.New(109021, "收款记录不存在")
	ReceiptRecordAlreadyExists = xcode.New(109022, "收款记录已存在")
	AddReceiptRecordFail       = xcode.New(109023, "添加收款记录失败")
	UpdateReceiptRecordFail    = xcode.New(109024, "更新收款记录失败")
	DeleteReceiptRecordFail    = xcode.New(109025, "删除收款记录失败")
	GetReceiptRecordFail       = xcode.New(109026, "获取收款记录失败")
	ReceiptAmountInvalid       = xcode.New(109027, "收款金额无效")

	// 工资发放相关错误 109031-109040
	SalaryPaymentNotFound      = xcode.New(109031, "工资发放记录不存在")
	SalaryPaymentAlreadyExists = xcode.New(109032, "工资发放记录已存在")
	AddSalaryPaymentFail       = xcode.New(109033, "添加工资发放记录失败")
	UpdateSalaryPaymentFail    = xcode.New(109034, "更新工资发放记录失败")
	DeleteSalaryPaymentFail    = xcode.New(109035, "删除工资发放记录失败")
	GetSalaryPaymentFail       = xcode.New(109036, "获取工资发放记录失败")
	SalaryAmountInvalid        = xcode.New(109037, "工资金额无效")
	PaymentStatusInvalid       = xcode.New(109038, "发放状态无效")
)

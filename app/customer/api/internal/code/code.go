package code

import "erp/common/xcode"

var (
	// 客户相关错误 103001-103010
	CustomerNotFound      = xcode.New(103001, "客户不存在")
	CustomerAlreadyExists = xcode.New(103002, "客户已存在")
	CustomerCodeDuplicate = xcode.New(103003, "客户编码重复")
	CustomerInUse         = xcode.New(103004, "客户正在使用中，无法删除")
	AddCustomerFail       = xcode.New(103005, "添加客户失败")
	UpdateCustomerFail    = xcode.New(103006, "更新客户失败")
	DeleteCustomerFail    = xcode.New(103007, "删除客户失败")
	CreditLimitExceeded   = xcode.New(103008, "客户信用额度已超限")

	// 客户分类相关错误 103011-103020
	CustomerCategoryNotFound      = xcode.New(103011, "客户分类不存在")
	CustomerCategoryAlreadyExists = xcode.New(103012, "客户分类已存在")
	AddCustomerCategoryFail       = xcode.New(103013, "添加客户分类失败")
	UpdateCustomerCategoryFail    = xcode.New(103014, "更新客户分类失败")
	DeleteCustomerCategoryFail    = xcode.New(103015, "删除客户分类失败")
	CustomerCategoryInUse         = xcode.New(103016, "客户分类正在使用中，无法删除")

	// 客户满意度调查相关错误 103021-103030
	CustomerSatisfactionNotFound = xcode.New(103021, "客户满意度调查不存在")
	AddSatisfactionSurveyFail    = xcode.New(103022, "添加客户满意度调查失败")
	GetSatisfactionSurveyFail    = xcode.New(103023, "获取客户满意度调查失败")
)

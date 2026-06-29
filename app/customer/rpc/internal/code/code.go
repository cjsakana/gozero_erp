package code

import "erp/common/xcode"

var (
	// 客户相关错误 103101-103110
	CustomerNotFound      = xcode.New(103101, "客户不存在")
	CustomerAlreadyExists = xcode.New(103102, "客户已存在")
	CustomerCodeDuplicate = xcode.New(103103, "客户编码重复")
	CustomerInUse         = xcode.New(103104, "客户正在使用中，无法删除")
	AddCustomerFail       = xcode.New(103105, "添加客户失败")
	UpdateCustomerFail    = xcode.New(103106, "更新客户失败")
	DeleteCustomerFail    = xcode.New(103107, "删除客户失败")
	CreditLimitExceeded   = xcode.New(103108, "客户信用额度已超限")
	GetCustomerFail       = xcode.New(103109, "获取客户信息失败")

	// 客户分类相关错误 103111-103120
	CustomerCategoryNotFound      = xcode.New(103111, "客户分类不存在")
	CustomerCategoryAlreadyExists = xcode.New(103112, "客户分类已存在")
	AddCustomerCategoryFail       = xcode.New(103113, "添加客户分类失败")
	UpdateCustomerCategoryFail    = xcode.New(103114, "更新客户分类失败")
	DeleteCustomerCategoryFail    = xcode.New(103115, "删除客户分类失败")
	CustomerCategoryInUse         = xcode.New(103116, "客户分类正在使用中，无法删除")
	GetCustomerCategoryFail       = xcode.New(103117, "获取客户分类失败")

	// 客户满意度调查相关错误 103121-103130
	CustomerSatisfactionNotFound = xcode.New(103121, "客户满意度调查不存在")
	AddSatisfactionSurveyFail    = xcode.New(103122, "添加客户满意度调查失败")
	GetSatisfactionSurveyFail    = xcode.New(103123, "获取客户满意度调查失败")
)

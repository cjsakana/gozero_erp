package code

import "erp/common/xcode"

var (
	// 供应商相关错误 104101-104110
	SupplierNotFound      = xcode.New(104101, "供应商不存在")
	SupplierAlreadyExists = xcode.New(104102, "供应商已存在")
	SupplierCodeDuplicate = xcode.New(104103, "供应商编码重复")
	SupplierInUse         = xcode.New(104104, "供应商正在使用中，无法删除")
	AddSupplierFail       = xcode.New(104105, "添加供应商失败")
	UpdateSupplierFail    = xcode.New(104106, "更新供应商失败")
	DeleteSupplierFail    = xcode.New(104107, "删除供应商失败")

	// 供应商评价相关错误 104111-104120
	SupplierEvaluationNotFound = xcode.New(104111, "供应商评价不存在")
	AddEvaluationFail          = xcode.New(104112, "添加供应商评价失败")
	GetEvaluationFail          = xcode.New(104113, "获取供应商评价失败")
	SearchEvaluationFail       = xcode.New(104114, "搜索供应商评价失败")
)

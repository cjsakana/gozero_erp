package code

import "erp/common/xcode"

var (
	// 供应商相关错误 104001-104010
	SupplierNotFound      = xcode.New(104001, "供应商不存在")
	SupplierAlreadyExists = xcode.New(104002, "供应商已存在")
	SupplierCodeDuplicate = xcode.New(104003, "供应商编码重复")
	SupplierInUse         = xcode.New(104004, "供应商正在使用中，无法删除")
	AddSupplierFail       = xcode.New(104005, "添加供应商失败")
	UpdateSupplierFail    = xcode.New(104006, "更新供应商失败")
	DeleteSupplierFail    = xcode.New(104007, "删除供应商失败")

	// 供应商评价相关错误 104011-104020
	SupplierEvaluationNotFound = xcode.New(104011, "供应商评价不存在")
	AddEvaluationFail          = xcode.New(104012, "添加供应商评价失败")
	GetEvaluationFail          = xcode.New(104013, "获取供应商评价失败")
	SearchEvaluationFail       = xcode.New(104014, "搜索供应商评价失败")
)

package code

import "erp/common/xcode"

var (
	// 产品相关错误 105001-105010
	ProductNotFound      = xcode.New(105001, "产品不存在")
	ProductAlreadyExists = xcode.New(105002, "产品已存在")
	ProductNoDuplicate   = xcode.New(105003, "产品编号重复")
	ProductInUse         = xcode.New(105004, "产品正在使用中，无法删除")
	AddProductFail       = xcode.New(105005, "添加产品失败")
	UpdateProductFail    = xcode.New(105006, "更新产品失败")
	DeleteProductFail    = xcode.New(105007, "删除产品失败")
	ImportProductFail    = xcode.New(105008, "导入产品失败")
	ExportProductFail    = xcode.New(105009, "导出产品失败")

	// 产品分类相关错误 105011-105020
	ProductCategoryNotFound      = xcode.New(105011, "产品分类不存在")
	ProductCategoryAlreadyExists = xcode.New(105012, "产品分类已存在")
	AddCategoryFail              = xcode.New(105013, "添加产品分类失败")
	UpdateCategoryFail           = xcode.New(105014, "更新产品分类失败")
	DeleteCategoryFail           = xcode.New(105015, "删除产品分类失败")
	CategoryInUse                = xcode.New(105016, "产品分类正在使用中，无法删除")

	// 产品批次相关错误 105021-105030
	ProductBatchNotFound = xcode.New(105021, "产品批次不存在")
	AddBatchFail         = xcode.New(105022, "添加产品批次失败")
	GetBatchFail         = xcode.New(105023, "获取产品批次失败")
	SearchBatchFail      = xcode.New(105024, "搜索产品批次失败")
)

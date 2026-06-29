package code

import "erp/common/xcode"

var (
	// 产品相关错误 105101-105110
	ProductNotFound      = xcode.New(105101, "产品不存在")
	ProductAlreadyExists = xcode.New(105102, "产品已存在")
	ProductNoDuplicate   = xcode.New(105103, "产品编号重复")
	ProductInUse         = xcode.New(105104, "产品正在使用中，无法删除")
	AddProductFail       = xcode.New(105105, "添加产品失败")
	UpdateProductFail    = xcode.New(105106, "更新产品失败")
	DeleteProductFail    = xcode.New(105107, "删除产品失败")
	ImportProductFail    = xcode.New(105108, "导入产品失败")
	GetProductFail       = xcode.New(105110, "获取产品信息失败")

	// 产品分类相关错误 105111-105120
	ProductCategoryNotFound      = xcode.New(105111, "产品分类不存在")
	ProductCategoryAlreadyExists = xcode.New(105112, "产品分类已存在")
	AddCategoryFail              = xcode.New(105113, "添加产品分类失败")
	UpdateCategoryFail           = xcode.New(105114, "更新产品分类失败")
	DeleteCategoryFail           = xcode.New(105115, "删除产品分类失败")
	CategoryInUse                = xcode.New(105116, "产品分类正在使用中，无法删除")
	GetProductCategoryFail       = xcode.New(105117, "获取产品分类信息失败")

	// 产品批次相关错误 105121-105130
	ProductBatchNotFound = xcode.New(105121, "产品批次不存在")
	AddBatchFail         = xcode.New(105122, "添加产品批次失败")
	GetBatchFail         = xcode.New(105123, "获取产品批次失败")
	SearchBatchFail      = xcode.New(105124, "搜索产品批次失败")
	GetProductBatchFail  = xcode.New(105125, "获取产品批次信息失败")
)

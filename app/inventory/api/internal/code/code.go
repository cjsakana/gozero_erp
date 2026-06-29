package code

import "erp/common/xcode"

var (
	// 库存相关错误 106001-106010
	InventoryNotFound      = xcode.New(106001, "库存不存在")
	InventoryAlreadyExists = xcode.New(106002, "库存已存在")
	AddInventoryFail       = xcode.New(106003, "添加库存失败")
	UpdateInventoryFail    = xcode.New(106004, "更新库存失败")
	StockInsufficient      = xcode.New(106005, "库存不足")
	StockLocked            = xcode.New(106006, "库存已锁定")
	AdjustInventoryFail    = xcode.New(106007, "调整库存失败")
	InventoryCheckFail     = xcode.New(106008, "库存盘点失败")
	GetInventoryFail       = xcode.New(106009, "获取库存信息失败")

	// 仓库相关错误 106011-106020
	WarehouseNotFound      = xcode.New(106011, "仓库不存在")
	WarehouseAlreadyExists = xcode.New(106012, "仓库已存在")
	WarehouseNoDuplicate   = xcode.New(106013, "仓库编号重复")
	WarehouseInUse         = xcode.New(106014, "仓库正在使用中，无法删除")
	AddWarehouseFail       = xcode.New(106015, "添加仓库失败")
	UpdateWarehouseFail    = xcode.New(106016, "更新仓库失败")
	DeleteWarehouseFail    = xcode.New(106017, "删除仓库失败")
	GetWarehouseFail       = xcode.New(106018, "获取仓库信息失败")

	// 库存交易相关错误 106021-106030
	InventoryTransactionNotFound = xcode.New(106021, "库存交易记录不存在")
	GetTransactionFail           = xcode.New(106022, "获取库存交易记录失败")
	CreateTransactionFail        = xcode.New(106023, "创建库存交易记录失败")

	ParamsInvalid = xcode.New(106024, "参数无效")
)

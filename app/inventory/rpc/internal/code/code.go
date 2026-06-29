package code

import "erp/common/xcode"

var (
	// 库存相关错误 106101-106110
	InventoryNotFound      = xcode.New(106101, "库存不存在")
	InventoryAlreadyExists = xcode.New(106102, "库存已存在")
	AddInventoryFail       = xcode.New(106103, "添加库存失败")
	UpdateInventoryFail    = xcode.New(106104, "更新库存失败")
	StockInsufficient      = xcode.New(106105, "库存不足")
	StockLocked            = xcode.New(106106, "库存已锁定")
	AdjustInventoryFail    = xcode.New(106107, "调整库存失败")
	InventoryCheckFail     = xcode.New(106108, "库存盘点失败")
	GetInventoryFail       = xcode.New(106109, "获取库存信息失败")

	// 仓库相关错误 106111-106120
	WarehouseNotFound      = xcode.New(106111, "仓库不存在")
	WarehouseAlreadyExists = xcode.New(106112, "仓库已存在")
	WarehouseNoDuplicate   = xcode.New(106113, "仓库编号重复")
	WarehouseInUse         = xcode.New(106114, "仓库正在使用中，无法删除")
	AddWarehouseFail       = xcode.New(106115, "添加仓库失败")
	UpdateWarehouseFail    = xcode.New(106116, "更新仓库失败")
	DeleteWarehouseFail    = xcode.New(106117, "删除仓库失败")
	GetWarehouseFail       = xcode.New(106118, "获取仓库信息失败")

	// 库存交易相关错误 106121-106130
	InventoryTransactionNotFound = xcode.New(106121, "库存交易记录不存在")
	GetTransactionFail           = xcode.New(106122, "获取库存交易记录失败")
	CreateTransactionFail        = xcode.New(106123, "创建库存交易记录失败")
)

package xtypes

// 权限常量定义
const (
	// 采购管理 (10xx)
	PurchaseOrderView    = 1001 // 查看采购单
	PurchaseOrderCreate  = 1002 // 创建采购单
	PurchaseOrderUpdate  = 1003 // 修改采购单
	PurchaseOrderSubmit  = 1004 // 提交采购单
	PurchaseOrderApprove = 1005 // 审批采购单
	PurchaseOrderCancel  = 1006 // 取消采购单
	PurchaseOrderReceive = 1007 // 采购收货
	SupplierView         = 1008 // 查看供应商
	SupplierManage       = 1009 // 管理供应商

	// 销售管理 (20xx)
	SalesOrderView    = 2001 // 查看销售单
	SalesOrderCreate  = 2002 // 创建销售单
	SalesOrderUpdate  = 2003 // 修改销售单
	SalesOrderSubmit  = 2004 // 提交销售单
	SalesOrderApprove = 2005 // 审批销售单
	SalesOrderCancel  = 2006 // 取消销售单
	SalesOrderDeliver = 2007 // 销售发货

	// 库存管理 (30xx)
	InventoryView     = 3001 // 查看库存
	InventoryAdjust   = 3002 // 库存调整
	InventoryTransfer = 3003 // 库存调拨
	InventoryCount    = 3004 // 库存盘点
	WarehouseManage   = 3005 // 管理仓库

	// 商品管理 (40xx)
	ProductView    = 4001 // 查看商品
	ProductManage  = 4002 // 管理商品
	CategoryManage = 4003 // 管理分类

	// 财务管理 (50xx)
	CashFlowView   = 5001 // 查看资金流水
	CashFlowManage = 5002 // 管理资金流水
	AccountManage  = 5003 // 管理银行账户

	// 人力资源 (60xx)
	EmployeeView            = 6001 // 查看员工
	EmployeeManage          = 6002 // 管理员工
	DepartmentManage        = 6003 // 管理部门
	PayrollManage           = 6004 // 管理薪酬
	EmployeeSelfManage      = 6005 // 管理个人信息
	EmployeeSensitiveManage = 6006 // 管理员工敏感信息

	// 客户关系 (70xx)
	CustomerView   = 7001 // 查看客户
	CustomerManage = 7002 // 管理客户
	SurveyManage   = 7003 // 管理满意度调查

	// 系统管理 (90xx)
	UserView       = 9001 // 查看用户
	UserManage     = 9002 // 管理用户
	RoleView       = 9003 // 查看角色
	RoleManage     = 9004 // 管理角色
	PermissionView = 9005 // 查看权限
	UserSelfManage = 9006 // 用户管理个人信息
)

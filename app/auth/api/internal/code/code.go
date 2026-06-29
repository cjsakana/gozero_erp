package code

import "erp/common/xcode"

var (
	// 角色相关错误 102001-102010
	RoleNotFound      = xcode.New(102001, "角色不存在")
	RoleAlreadyExists = xcode.New(102002, "角色已存在")
	RoleCodeDuplicate = xcode.New(102003, "角色编码重复")
	RoleInUse         = xcode.New(102004, "角色正在使用中，无法删除")
	AddRoleFail       = xcode.New(102005, "添加角色失败")
	UpdateRoleFail    = xcode.New(102006, "更新角色失败")
	DeleteRoleFail    = xcode.New(102007, "删除角色失败")

	// 权限相关错误 102011-102020
	PermissionNotFound = xcode.New(102011, "权限不存在")
	GetPermissionFail  = xcode.New(102012, "获取权限列表失败")

	// 角色权限相关错误 102021-102030
	RolePermissionNotFound      = xcode.New(102021, "角色权限关系不存在")
	RolePermissionAlreadyExists = xcode.New(102022, "角色权限关系已存在")
	AddRolePermissionFail       = xcode.New(102023, "添加角色权限失败")
	UpdateRolePermissionFail    = xcode.New(102024, "更新角色权限失败")
	DeleteRolePermissionFail    = xcode.New(102025, "删除角色权限失败")

	// 用户角色相关错误 102031-102040
	UserRoleNotFound      = xcode.New(102031, "用户角色关系不存在")
	UserRoleAlreadyExists = xcode.New(102032, "用户角色关系已存在")
	AddUserRoleFail       = xcode.New(102033, "添加用户角色失败")
	UpdateUserRoleFail    = xcode.New(102034, "更新用户角色失败")
	DeleteUserRoleFail    = xcode.New(102035, "删除用户角色失败")
)

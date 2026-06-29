package code

import "erp/common/xcode"

var (
	// 角色相关错误 102101-102110
	RoleNotFound      = xcode.New(102101, "角色不存在")
	RoleAlreadyExists = xcode.New(102102, "角色已存在")
	RoleCodeDuplicate = xcode.New(102103, "角色编码重复")
	RoleInUse         = xcode.New(102104, "角色正在使用中，无法删除")
	AddRoleFail       = xcode.New(102105, "添加角色失败")
	UpdateRoleFail    = xcode.New(102106, "更新角色失败")
	DeleteRoleFail    = xcode.New(102107, "删除角色失败")

	// 权限相关错误 102111-102120
	PermissionNotFound = xcode.New(102111, "权限不存在")
	GetPermissionFail  = xcode.New(102112, "获取权限列表失败")

	// 角色权限相关错误 102121-102130
	RolePermissionNotFound      = xcode.New(102121, "角色权限关系不存在")
	RolePermissionAlreadyExists = xcode.New(102122, "角色权限关系已存在")
	AddRolePermissionFail       = xcode.New(102123, "添加角色权限失败")
	UpdateRolePermissionFail    = xcode.New(102124, "更新角色权限失败")
	DeleteRolePermissionFail    = xcode.New(102125, "删除角色权限失败")

	// 用户角色相关错误 102131-102140
	UserRoleNotFound      = xcode.New(102131, "用户角色关系不存在")
	UserRoleAlreadyExists = xcode.New(102132, "用户角色关系已存在")
	AddUserRoleFail       = xcode.New(102133, "添加用户角色失败")
	UpdateUserRoleFail    = xcode.New(102134, "更新用户角色失败")
	DeleteUserRoleFail    = xcode.New(102135, "删除用户角色失败")
)

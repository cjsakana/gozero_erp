package code

import "erp/common/xcode"

var (
	// 登录认证相关错误 101101-101110
	LoginFail                = xcode.New(101101, "账户或密码错误")
	VerificationCodeInvalid  = xcode.New(101102, "验证码无效")
	SendVerificationCodeFail = xcode.New(101103, "发送验证码失败，请重试")
	TokenInvalid             = xcode.New(101104, "Token无效")
	TokenExpired             = xcode.New(101105, "Token已过期")

	// 密码相关错误 101111-101120
	SamePasswordTwice    = xcode.New(101111, "新旧密码相同")
	ForgotPasswordFail   = xcode.New(101112, "忘记密码失败，请重试")
	UpdatePasswordFail   = xcode.New(101113, "更新密码失败")
	OldPasswordIncorrect = xcode.New(101114, "原密码错误")
	PasswordTooWeak      = xcode.New(101115, "密码强度不够")

	// 用户信息相关错误 101121-101130
	UserNotFound        = xcode.New(101121, "用户不存在")
	UserAlreadyExists   = xcode.New(101122, "用户已存在")
	UsernameDuplicate   = xcode.New(101123, "用户名重复")
	PhoneDuplicate      = xcode.New(101124, "手机号重复")
	EmployeeNoDuplicate = xcode.New(101125, "员工编号重复")
	AddUserFail         = xcode.New(101126, "添加用户失败")
	UpdateUserFail      = xcode.New(101127, "更新用户失败")
	DeleteUserFail      = xcode.New(101128, "删除用户失败")
	GetUserFail         = xcode.New(101129, "获取用户信息失败")
	SearchUserFail      = xcode.New(101130, "搜索用户失败")
	BulkInsertUserFail  = xcode.New(101131, "批量插入用户失败")
)

package code

import "erp/common/xcode"

var (
	// 文件上传相关错误 101001-101010
	PutBucketErr     = xcode.New(101001, "上传bucket失败")
	UploadAvatarFail = xcode.New(101007, "上传头像失败")

	// 登录认证相关错误 101011-101020
	LoginFail                = xcode.New(101002, "账户或密码错误")
	VerificationCodeInvalid  = xcode.New(101003, "验证码无效")
	SendVerificationCodeFail = xcode.New(101005, "发送验证码失败，请重试")
	TokenInvalid             = xcode.New(101008, "Token无效")
	TokenExpired             = xcode.New(101009, "Token已过期")
	RefreshTokenFail         = xcode.New(101010, "刷新Token失败")
	LogoutFail               = xcode.New(101011, "退出登录失败")

	// 密码相关错误 101021-101030
	SamePasswordTwice    = xcode.New(101004, "新旧密码相同")
	ForgotPasswordFail   = xcode.New(101006, "忘记密码失败，请重试")
	UpdatePasswordFail   = xcode.New(101012, "更新密码失败")
	OldPasswordIncorrect = xcode.New(101013, "原密码错误")
	PasswordTooWeak      = xcode.New(101014, "密码强度不够")

	// 用户信息相关错误 101031-101040
	UserNotFound       = xcode.New(101031, "用户不存在")
	UserAlreadyExists  = xcode.New(101032, "用户已存在")
	UsernameDuplicate  = xcode.New(101033, "用户名重复")
	PhoneDuplicate     = xcode.New(101034, "手机号重复")
	GetUserFail        = xcode.New(101035, "获取用户信息失败")
	UpdateUserInfoFail = xcode.New(101036, "更新用户信息失败")
	SearchUserFail     = xcode.New(101037, "搜索用户失败")
)

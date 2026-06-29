package types

type (
	SearchUserParams struct {
		Page       int64
		Limit      int64
		EmployeeId int64
		Username   string
		RealName   string
		Phone      string
		Email      string
		Resigned   *bool // 指针类型，允许为 nil 表示不筛选
	}
)

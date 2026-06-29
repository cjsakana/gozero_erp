package types

type SearchUserRole struct {
	Page       int64
	Limit      int64
	EmployeeNo string
	RealName   string
	RoleId     int64
}

type SearchUserRoleItem struct {
	Id          int64
	EmployeeNo  string
	RealName    string
	Code        string
	Name        string
	Description string
	RoleId      int64
}

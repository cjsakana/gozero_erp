package types

type (
	SearchDepartmentParams struct {
		SearchCom
		Name     string
		ParentId int64
		Code     string
	}
)

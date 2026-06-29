package types

import "time"

type (
	SearchEmployeeDetailParam struct {
		SearchCom
		Gender       int64
		DepartmentId int64
		PositionId   int64
		Salary       float64
		HireDate     time.Time
		Name         string
		Resigned     int64
	}
)

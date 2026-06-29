package types

type (
	SearchProductCategoryParams struct {
		SearchCom
		CategoryName string
		ParentId     int64
	}
)

package types

type (
	SearchProductParams struct {
		SearchCom
		ProductNo   string
		ProductName string
		CategoryId  int64
		IsActive    int64
		IsMaterial  int64
	}
)

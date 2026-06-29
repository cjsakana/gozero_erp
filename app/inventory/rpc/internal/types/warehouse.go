package types

type (
	SearchWarehouseParams struct {
		SearchCom
		Name     string
		Location string
		IsActive int64
	}
)

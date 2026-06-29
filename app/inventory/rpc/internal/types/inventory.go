package types

type (
	SearchInventoryParams struct {
		SearchCom
		ProductId    int64
		WarehouseId  int64
		CurrentStock float64
		SafetyStock  float64
		LockedStock  float64
	}
)

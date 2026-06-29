package types

type GetRequisitionListParams struct {
	Page        int64
	PageSize    int64
	WorkOrderId int64
	WarehouseId int64
	Status      int64
}

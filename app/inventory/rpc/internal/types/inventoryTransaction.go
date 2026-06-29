package types

import "time"

type (
	SearchInventoryTransactionParams struct {
		SearchCom
		ProductId       int64
		WarehouseId     int64
		BatchId         int64
		TransactionType int64
		ReferenceType   int64
		StartTime       time.Time
		EndTime         time.Time
	}
)

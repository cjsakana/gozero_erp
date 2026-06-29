package types

import "time"

type (
	SearchProductBatchParams struct {
		SearchCom
		ProductId int64
		BatchNo   string
		StartDate time.Time
		EndDate   time.Time
	}
)

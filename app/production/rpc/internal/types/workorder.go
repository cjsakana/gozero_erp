package types

import "time"

type GetWorkOrderListParams struct {
	Page      int64
	PageSize  int64
	ProductId int64
	Status    int64
	Priority  int64
	StartDate time.Time
	EndDate   time.Time
}

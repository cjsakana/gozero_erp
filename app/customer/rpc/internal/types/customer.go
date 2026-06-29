package types

type SearchCustomer struct {
	SearchCom
	Code         string
	USCC         string
	Name         string
	CategoryId   int64
	Contact      string
	Address      string
	PaymentTerms string
	IsActive     int64
}

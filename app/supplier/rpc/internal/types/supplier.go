package types

type (
	SearchSupplierParams struct {
		SearchCom
		Code         string
		Uscc         string
		Name         string
		Contact      string
		Address      string
		PaymentTerms string
		Credit       string
		IsActive     int64
	}
)

package types

type (
	SearchSupplierEvaluation struct {
		SearchCom
		SupplierId  int64
		QualityMin  float64
		QualityMax  float64
		QualityOp   string //质量评分操作符: gt/gte/eq/lt/lte
		DeliveryMin float64
		DeliveryMax float64
		DeliveryOp  string
		ServiceMin  float64
		ServiceMax  float64
		ServiceOp   string
		OverallMin  float64
		OverallMax  float64
		OverallOp   string
		StartData   int64
		EndData     int64
		EvaluatorId int64
	}
)

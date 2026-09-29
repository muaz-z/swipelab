package merchant

type Merchant struct {
	ID     string
	Name   string
	Active bool
}

var merchants = map[string]Merchant{
	"MERCHANT_001": {
		ID:     "MERCHANT_001",
		Name:   "SwipeMart",
		Active: true,
	},
	"MERCHANT_002": {
		ID:     "MERCHANT_002",
		Name:   "Suspended Shop",
		Active: false,
	},
}

func FindByID(id string) (Merchant, bool) {
	merchant, exists := merchants[id]

	return merchant, exists
}

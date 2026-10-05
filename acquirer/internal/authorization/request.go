package authorization

type Request struct {
	MerchantID string `json:"merchant_id"`
	CardNumber string `json:"card_number"`
	Amount     int64  `json:"amount"`
	Currency   string `json:"currency"`
}

package network

type AuthorizationRequest struct {
	AuthorizationID string `json:"authorization_id"`
	MerchantID      string `json:"merchant_id"`
	CardNumber      string `json:"card_number"`
	Amount          int64  `json:"amount"`
	Currency        string `json:"currency"`
}

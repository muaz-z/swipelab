package authorization

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending  Status = "PENDING"
	StatusApproved Status = "APPROVED"
	StatusDeclined Status = "DECLINED"
	StatusFailed   Status = "FAILED"
)

type Authorization struct {
	ID                 string    `json:"id"`
	MerchantID         string    `json:"merchant_id"`
	Amount             int64     `json:"amount"`
	Currency           string    `json:"currency"`
	Status             Status    `json:"status"`
	IdempotencyKey     string    `json:"-"`
	RequestFingerprint string    `json:"-"`
	CreatedAt          time.Time `json:"created_at"`
}

func New(request Request, idempotencyKey string, requestFingerPrint string) Authorization {
	return Authorization{
		ID:                 uuid.NewString(),
		MerchantID:         request.MerchantID,
		Amount:             request.Amount,
		Currency:           request.Currency,
		Status:             StatusPending,
		IdempotencyKey:     idempotencyKey,
		RequestFingerprint: requestFingerPrint,
		CreatedAt:          time.Now().UTC(),
	}
}

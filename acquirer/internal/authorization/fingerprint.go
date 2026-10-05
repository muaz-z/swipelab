package authorization

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func Fingerprint(request Request) string {
	value := fmt.Sprintf(
		"%s|%s|%d|%s",
		request.MerchantID,
		request.CardNumber,
		request.Amount,
		request.Currency,
	)

	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}

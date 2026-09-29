package authorization

import "fmt"

func Validate(request Request) error {
	if request.MerchantID == "" {
		return fmt.Errorf("merchant_id is required")
	}

	if request.Amount <= 0 {
		return fmt.Errorf("amount must be bigger than 0")
	}

	switch request.Currency {
	case "HUF", "EUR", "USD":
		// VALID
	default:
		return fmt.Errorf("unsupported currency")
	}

	return nil
}

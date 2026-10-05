package authorization

import "fmt"

func Validate(request Request) error {
	if request.AuthorizationID == "" {
		return fmt.Errorf("authorization_id is required")
	}

	if request.MerchantID == "" {
		return fmt.Errorf("merchant_id is required")
	}

	if request.CardNumber == "" {
		return fmt.Errorf("card_number is required")
	}

	if request.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	switch request.Currency {
	case "HUF", "EUR", "USD":

	default:
		return fmt.Errorf("unsupported currency")
	}

	return nil
}

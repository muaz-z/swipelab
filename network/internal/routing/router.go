package routing

import (
	"fmt"
	"strings"
)

type Issuer struct {
	Name string
	URL  string
}

func FindIssuer(cardNumber string) (Issuer, error) {
	switch {
	case strings.HasPrefix(cardNumber, "411111"):
		return Issuer{
			Name: "SwipeBank",
			URL:  "http://localhost:8083",
		}, nil
	case strings.HasPrefix(cardNumber, "550000"):
		return Issuer{
			Name: "DemoBank",
			URL:  "http://localhost:8084",
		}, nil
	default:
		return Issuer{}, fmt.Errorf("issuer not found")
	}

}

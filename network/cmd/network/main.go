package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/muaz-z/swipelab/network/internal/authorization"
	"github.com/muaz-z/swipelab/network/internal/issuer"
	"github.com/muaz-z/swipelab/network/internal/routing"
)

func main() {
	issuerClient := issuer.NewClient()

	http.HandleFunc("/authorizations", func(w http.ResponseWriter, r *http.Request) {
		handleAuthorization(w, r, issuerClient)
	})
	fmt.Println("Card Network listening on :8082")

	err := http.ListenAndServe(":8082", nil)

	if err != nil {
		panic(err)
	}

}

func handleAuthorization(w http.ResponseWriter, r *http.Request, issuerClient *issuer.Client) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request authorization.Request

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = authorization.Validate(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	issuerInfo, err := routing.FindIssuer(request.CardNumber)
	if err != nil {
		http.Error(w, "issuer not found", http.StatusBadRequest)
		return
	}

	err = issuerClient.Authorize(r.Context(), issuerInfo.URL, request)
	if err != nil {
		fmt.Printf("issuer authorization failed: %v/n", err)
		http.Error(w, "issuer processing failed", http.StatusBadGateway)
		return
	}

	fmt.Printf(
		"Authorization received id=%s merchant=%s amount=%d currency=%s\n",
		request.AuthorizationID,
		request.MerchantID,
		request.Amount,
		request.Currency,
	)

	fmt.Printf(
		"Authorization %s routed to %s\n",
		request.AuthorizationID,
		issuerInfo.Name,
	)

	w.WriteHeader(http.StatusOK)
}

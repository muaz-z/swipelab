package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/muaz-z/swipelab/acquirer/internal/authorization"
	"github.com/muaz-z/swipelab/acquirer/internal/database"
	"github.com/muaz-z/swipelab/acquirer/internal/merchant"
	"github.com/muaz-z/swipelab/acquirer/internal/network"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("No .env file loaded")
	}
	ctx := context.Background()

	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		panic("DATABASE_URL is not set")
	}

	db, err := database.Connect(ctx, databaseUrl)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	authorizationRepository := authorization.NewRepository(db)
	fmt.Println("Connected to POSTGRESQL")

	networkClient := network.NewClient("http://localhost:8082")

	http.HandleFunc("/authorizations", func(w http.ResponseWriter, r *http.Request) {
		handleAuthorization(w, r, authorizationRepository, networkClient)
	})
	fmt.Println("Acquirer listening on :8081")

	err = http.ListenAndServe(":8081", nil)

	if err != nil {
		panic(err)
	}
}

func handleAuthorization(w http.ResponseWriter, r *http.Request, repository *authorization.Repository, networkClient *network.Client) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowerd", http.StatusMethodNotAllowed)
		return
	}

	idempotencyKey := r.Header.Get("Idempotency-key")

	if idempotencyKey == "" {
		http.Error(
			w,
			"Idempotency-Key header is required",
			http.StatusBadRequest,
		)
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

	requestFingerprint := authorization.Fingerprint(request)

	m, exists := merchant.FindByID(request.MerchantID)

	if !exists {
		http.Error(w, "merchant not found", http.StatusBadRequest)
		return

	}

	if !m.Active {
		http.Error(w, "merchant is inactive", http.StatusForbidden)
		return

	}

	existingAuth, err := repository.FindByIdempotencyKey(
		r.Context(),
		request.MerchantID,
		idempotencyKey,
	)

	if err != nil {
		fmt.Println("failed to check idempotency key", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if existingAuth != nil {
		if existingAuth.RequestFingerprint != requestFingerprint {
			http.Error(
				w,
				"idempotency key already used with a different request",
				http.StatusConflict,
			)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(existingAuth)
		return
	}

	auth := authorization.New(request, idempotencyKey, requestFingerprint)

	err = repository.Create(r.Context(), auth)

	if errors.Is(err, authorization.ErrDuplicateIdempotencyKey) {
		existingAuth, findErr := repository.FindByIdempotencyKey(
			r.Context(),
			request.MerchantID,
			idempotencyKey,
		)

		if findErr != nil {
			fmt.Println("failed to find existing authorization:", findErr)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if existingAuth == nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if existingAuth.RequestFingerprint != requestFingerprint {
			http.Error(
				w,
				"idempotency key already used with a different request",
				http.StatusConflict,
			)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(existingAuth)
		return
	}

	if err != nil {
		fmt.Println("failed to create authorization:", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	networkRequest := network.AuthorizationRequest{
		AuthorizationID: auth.ID,
		MerchantID:      auth.MerchantID,
		CardNumber:      request.CardNumber,
		Amount:          auth.Amount,
		Currency:        auth.Currency,
	}

	err = networkClient.Authorize(r.Context(), networkRequest)

	if err != nil {
		fmt.Println("failed to authorize card network:", err)

		if errors.Is(err, network.ErrRejected) {
			updateErr := repository.UpdateStatus(r.Context(), auth.ID, authorization.StatusFailed)
			if updateErr != nil {
				fmt.Println("failed to update authorization status", updateErr)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			http.Error(w, "authorization processing failed", http.StatusBadGateway)
			return
		}
		http.Error(w, "authorization outcome unknown", http.StatusBadGateway)
		return

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(auth)

	fmt.Printf(
		"Created authorization: id=%s merchant=%s amount=%d currency=%s status=%s\n",
		auth.ID,
		auth.MerchantID,
		auth.Amount,
		auth.Currency,
		auth.Status,
	)
}

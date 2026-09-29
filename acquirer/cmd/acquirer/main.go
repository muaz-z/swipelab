package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/muaz-z/swipelab/acquirer/internal/authorization"
	"github.com/muaz-z/swipelab/acquirer/internal/database"
	"github.com/muaz-z/swipelab/acquirer/internal/merchant"
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

	http.HandleFunc("/authorizations", func(w http.ResponseWriter, r *http.Request) {
		handleAuthorization(w, r, authorizationRepository)
	})
	fmt.Println("Acquirer listening on :8081")

	err = http.ListenAndServe(":8081", nil)

	if err != nil {
		panic(err)
	}
}

func handleAuthorization(w http.ResponseWriter, r *http.Request, repository *authorization.Repository) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowerd", http.StatusMethodNotAllowed)
		return
	}

	idempotencyKey := r.Header.Get("idempotency-key")

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

	m, exists := merchant.FindByID(request.MerchantID)

	if !exists {
		http.Error(w, "merchant not found", http.StatusBadRequest)
		return

	}

	if !m.Active {
		http.Error(w, "merchant is inactive", http.StatusBadRequest)
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
	}

	if existingAuth != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(existingAuth)
		return
	}

	auth := authorization.New(request, idempotencyKey)

	err = repository.Create(r.Context(), auth)
	if err != nil {
		fmt.Println("failed to create authorization:", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
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

package authorization

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

var ErrDuplicateIdempotencyKey = errors.New("duplicate idempotency key")

func (r *Repository) Create(ctx context.Context, auth Authorization) error {
	query := `
		INSERT INTO authorizations (
			id,
			merchant_id,
			amount,
			currency,
			status,
			idempotency_key,
			request_fingerprint,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		auth.ID,
		auth.MerchantID,
		auth.Amount,
		auth.Currency,
		auth.Status,
		auth.IdempotencyKey,
		auth.RequestFingerprint,
		auth.CreatedAt,
		auth.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrDuplicateIdempotencyKey
		}
		return err
	}

	return nil
}

func (r *Repository) FindByIdempotencyKey(
	ctx context.Context,
	merchantID string,
	idempotencyKey string,
) (*Authorization, error) {
	query := `
		SELECT
			id,
			merchant_id,
			amount,
			currency,
			status,
			idempotency_key,
			request_fingerprint,
			created_at
		FROM authorizations
		WHERE merchant_id = $1
		  AND idempotency_key = $2
	`

	var auth Authorization

	err := r.db.QueryRow(
		ctx,
		query,
		merchantID,
		idempotencyKey,
	).Scan(
		&auth.ID,
		&auth.MerchantID,
		&auth.Amount,
		&auth.Currency,
		&auth.Status,
		&auth.IdempotencyKey,
		&auth.RequestFingerprint,
		&auth.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &auth, nil

}

func (r *Repository) UpdateStatus(
	ctx context.Context,
	id string,
	status Status,
) error {
	query := `
	UPDATE authorizations
        SET status = $1,
            updated_at = NOW()
        WHERE id = $2
	`

	result, err := r.db.Exec(ctx, query, status, id)

	if err != nil {
		return fmt.Errorf("update authorization status: %w", err)

	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("authorization %s not found", id)
	}

	return nil
}

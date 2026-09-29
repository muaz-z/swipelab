package authorization

import (
	"context"

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

func (r *Repository) Create(ctx context.Context, auth Authorization) error {
	query := `
		INSERT INTO authorizations (
			id,
			merchant_id,
			amount,
			currency,
			status,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		auth.ID,
		auth.MerchantID,
		auth.Amount,
		auth.Currency,
		auth.Status,
		auth.CreatedAt,
		auth.CreatedAt,
	)

	return err
}

package repo

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/ozonmp/cnm-purchase-api/internal/model"
)

// Repo is DAO for Purchase
type Repo interface {
	DescribePurchase(ctx context.Context, purchaseID uint64) (*model.Purchase, error)
}

type repo struct {
	db        *sqlx.DB
	batchSize uint
}

// NewRepo returns Repo interface
func NewRepo(db *sqlx.DB, batchSize uint) Repo {
	return &repo{db: db, batchSize: batchSize}
}

func (r *repo) DescribePurchase(ctx context.Context, purchaseID uint64) (*model.Purchase, error) {
	return nil, nil
}

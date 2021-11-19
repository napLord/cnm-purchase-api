package repo

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/ozonmp/cnm-purchase-api/internal/model"
)

// Repo is DAO for Purchase
type Repo interface {
	DescribePurchase(ctx context.Context, purchaseID uint64) (*model.Purchase, error)
	CreatePurchase(ctx context.Context, totalSum uint64) (purchaseID uint64, err error)
	RemovePurchase(ctx context.Context, purchaseID uint64) (found bool, err error)
	ListPurchases(ctx context.Context, limit uint64, cursor uint64) (purchases []model.Purchase, nextCursor uint64, err error)
}

type repo struct {
	db        *sqlx.DB
	batchSize uint
}

var ErrNoPurchases = errors.New("got 0 purchases from db")

// NewRepo returns Repo interface
func NewRepo(db *sqlx.DB, batchSize uint) Repo {
	return &repo{db: db, batchSize: batchSize}
}

func (r *repo) DescribePurchase(ctx context.Context, purchaseID uint64) (*model.Purchase, error) {
	req := sq.
		Select("id", "total_sum").
		From("purchases").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": purchaseID}).
		Where(sq.Eq{"removed": false})

	s, args, err := req.ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "can't create sql query")
	}

	var res []model.Purchase

	err = r.db.SelectContext(ctx, &res, s, args...)

	if err != nil {
		return nil, errors.Wrap(err, "can't get purchase from db")
	}

	if len(res) == 0 {
		return nil, ErrNoPurchases
	}

	return &res[0], nil
}

func (r *repo) CreatePurchase(ctx context.Context, totalSum uint64) (purchaseID uint64, err error) {
	//transaction?
	req := sq.
		Insert("purchases").
		PlaceholderFormat(sq.Dollar).
		Columns("total_sum", "removed", "created").
		Suffix("RETURNING id").
		Values(totalSum, false, time.Now()).RunWith(r.db)

	rows, err := req.QueryContext(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "can't create purchase in db")
	}

	if rows.Next() {
		var lastInsertedID uint64
		err = rows.Scan(&lastInsertedID)
		if err != nil {
			return 0, errors.Wrap(err, "error during parsing row")
		}

		return uint64(lastInsertedID), nil
	}

	return 0, errors.New("qurety returned no rows")
}

func (r *repo) RemovePurchase(ctx context.Context, purchaseID uint64) (found bool, err error) {
	req := sq.
		Update("purchases").
		Set("removed", true).
		Set("updated", time.Now()).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": []uint64{purchaseID}}).
		RunWith(r.db)

	res, err := req.ExecContext(ctx)
	if err != nil {
		return false, errors.Wrap(err, "can't set purchase deleted in db")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, errors.Wrap(err, "can't get rowsAffected")
	}

	if rowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (r *repo) ListPurchases(ctx context.Context, limit uint64, cursor uint64) (purchases []model.Purchase, nextCursor uint64, err error) {
	req := sq.
		Select("id", "total_sum").
		From("purchases").
		PlaceholderFormat(sq.Dollar).
		Where(sq.GtOrEq{"id": cursor}).
		Where(sq.Eq{"removed": false}).
		OrderBy("id").
		Limit(limit)

	s, args, err := req.ToSql()

	if err != nil {
		return nil, 0, errors.Wrap(err, "can't create sql query")
	}

	var res []model.Purchase

	err = r.db.SelectContext(ctx, &res, s, args...)

	if err != nil {
		return nil, 0, errors.Wrap(err, "can't get purchases from db")
	}

	if len(res) == 0 {
		return nil, 0, ErrNoPurchases
	}

	return res, res[len(res)-1].ID + 1, nil
}

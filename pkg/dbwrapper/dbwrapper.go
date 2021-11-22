package dbwrapper

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

//WithTxFunc is a function which wrap a transaction actions
type WithTxFunc func(ctx context.Context, tx *sqlx.Tx) error

//WithTx wrapiing a transaction actions and rollback them if error occured or commit them otherwise
func WithTx(ctx context.Context, db *sqlx.DB, fn WithTxFunc) error {
	t, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "db.BeginTxx()")
	}

	if err = fn(ctx, t); err != nil {
		if errRollback := t.Rollback(); errRollback != nil {
			return errors.Wrap(err, "Tx.Rollback")
		}
		return errors.Wrap(err, "Tx.WithTxFunc")
	}

	if err = t.Commit(); err != nil {
		return errors.Wrap(err, "Tx.Commit")
	}
	return nil
}

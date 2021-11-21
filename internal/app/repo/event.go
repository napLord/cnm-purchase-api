package repo

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/ozonmp/cnm-purchase-api/internal/model"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
)

//EventRepo is interface of db
type EventRepo interface {
	Lock(ctx context.Context, n uint64) ([]model.PurchaseEvent, error)
	Unlock(ctx context.Context, eventIDs []uint64) error

	Add(ctx context.Context, event *model.PurchaseEvent) error
	Remove(ctx context.Context, eventIDs []uint64) error
}

type EventRepoImpl struct {
	db *sqlx.DB
}

func NewEventRepo(db *sqlx.DB) EventRepo {
	return &EventRepoImpl{db: db}
}

func (r *EventRepoImpl) Lock(ctx context.Context, n uint64) ([]model.PurchaseEvent, error) {
	selectVacantReq := sq.
		Select("id").
		PlaceholderFormat(sq.Dollar).
		From("purchases_events").
		Where(sq.Eq{"status": model.Unlocked.String()}).
		Limit(n)

	lockReq := sq.
		Update("purchases_events").
		PlaceholderFormat(sq.Dollar).
		Set("status", model.Locked.String()).
		Where(selectVacantReq.Prefix("id IN (").Suffix(")")).
		Suffix("RETURNING id, type, status, payload").
		RunWith(r.db)

	s, args, err := lockReq.ToSql()
	if err != nil {
		return nil, err
	}

	var result []model.PurchaseEvent
	err = r.db.SelectContext(ctx, &result, s, args...)

	return result, err
}

func (r *EventRepoImpl) Unlock(ctx context.Context, eventIDs []uint64) error {
	req := sq.
		Update("purchases_events").
		PlaceholderFormat(sq.Dollar).
		Set("status", model.Unlocked.String()).
		Where(sq.Eq{"id": eventIDs}).
		RunWith(r.db)

	_, err := req.ExecContext(ctx)
	if err != nil {
		return errors.New("can't unlock events in db")
	}

	return nil
}

func (r *EventRepoImpl) Add(ctx context.Context, event *model.PurchaseEvent) error {
	innerPurchasePB := model.PurchaseToPB(event.Entity)

	srlzdPurchase, err := protojson.Marshal(&innerPurchasePB)
	if err != nil {
		return errors.Wrap(err, "can't serialize purchase from event")
	}

	query := sq.
		Insert("purchases_events").
		PlaceholderFormat(sq.Dollar).
		Columns("purchase_id", "type", "status", "payload", "updated").
		Values(event.Entity.ID, event.Type.String(), model.Unlocked.String(), srlzdPurchase, time.Now()).
		RunWith(r.db)

	_, err = query.ExecContext(ctx)
	if err != nil {
		return errors.Wrap(err, "can't insert into db")
	}

	return nil
}

func (r *EventRepoImpl) Remove(ctx context.Context, eventIDs []uint64) error {
	query := sq.
		Delete("purchases_events").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": eventIDs}).
		RunWith(r.db)

	_, err := query.ExecContext(ctx)
	if err != nil {
		return errors.Wrap(err, "can't remove events from db")
	}

	return nil
}

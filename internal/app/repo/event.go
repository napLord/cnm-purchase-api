package repo

import (
	"github.com/napLord/cnm-purchase-api/internal/model"
)

//EventRepo is interface of db
type EventRepo interface {
	Lock(n uint64) ([]model.PurchaseEvent, error)
	Unlock(eventIDs []uint64) error

	Add(event []model.PurchaseEvent) error
	Remove(eventIDs []uint64) error
}

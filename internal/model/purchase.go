package model

// Purchase - purchase entity.
type Purchase struct {
	ID       uint64 `db:"id"`
	TotalSum uint64 `db:"total_sum"`
}

//EventType is type of Event
type EventType uint8

//EventStatus is status of event
type EventStatus uint8

//EventTypes
const (
	Created EventType = iota
	Updated
	Removed

	Deferred EventStatus = iota
	Processed
)

//PurchaseEvent is an event of purchasing something
type PurchaseEvent struct {
	ID     uint64
	Type   EventType
	Status EventStatus
	Entity *Purchase
}

//NewPurchaseEvent creates new PurchaseEvent
func NewPurchaseEvent(
	ID uint64,
	Entity *Purchase,
) *PurchaseEvent {
	return &PurchaseEvent{
		ID:     ID,
		Entity: Entity,
	}
}

package model

type Purchase struct {
	ID         uint64
	totalPrice uint64
}

type EventType uint8

type EventStatus uint8

const (
	Created EventType = iota
	Updated
	Removed

	Deferred EventStatus = iota
	Processed
)

type PurchaseEvent struct {
	ID     uint64
	Type   EventType
	Status EventStatus
	Entity *Purchase
}

func newPurchaseEvent(
	ID uint64,
	Type EventType,
	Status EventStatus,
	Entity *Purchase,
) *PurchaseEvent {
	return &PurchaseEvent{
		ID:     ID,
		Type:   Type,
		Status: Status,
		Entity: Entity,
	}
}

package model

type Purchase struct {
	ID         uint64
	TotalPrice uint64
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

func NewPurchaseEvent(
	ID uint64,
	Entity *Purchase,
) *PurchaseEvent {
	return &PurchaseEvent{
		ID:     ID,
		Entity: Entity,
	}
}

package model

import (
	"github.com/pkg/errors"

	pb "github.com/ozonmp/cnm-purchase-api/pkg/cnm-purchase-api"
	"google.golang.org/protobuf/encoding/protojson"
)

// Purchase - purchase entity.
type Purchase struct {
	ID       uint64 `db:"id"`
	TotalSum uint64 `db:"total_sum"`
}

func (t *Purchase) Scan(v interface{}) error {
	typeValue, ok := v.([]byte)
	if !ok {
		return errors.New("Scan error. interface{} is not a []byte")
	}

	var PBPurchae pb.Purchase

	err := protojson.Unmarshal(typeValue, &PBPurchae)
	if err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	*t = PBToPurchase(&PBPurchae)

	return nil
}

//EventType is type of Event
type EventType uint8

func (t EventType) String() string {
	switch t {
	case Created:
		return "Created"
	case Updated:
		return "Updated"
	case Removed:
		return "Removed"
	default:
		panic("unknow EventType")
	}
}

func (t *EventType) Scan(v interface{}) error {
	typeValue, ok := v.(string)
	if !ok {
		return errors.New("EventType is not a string somehow")
	}

	switch typeValue {
	case Created.String():
		*t = Created
	case Updated.String():
		*t = Updated
	case Removed.String():
		*t = Removed
	default:
		panic("unknow EventType")
	}

	return nil
}

//EventStatus is status of event
type EventStatus uint8

func (t EventStatus) String() string {
	switch t {
	case Locked:
		return "locked"
	case Unlocked:
		return "unlocked"
	default:
		panic("unknow EventStatus")
	}
}

func (t *EventStatus) Scan(v interface{}) error {
	typeValue, ok := v.(string)
	if !ok {
		return errors.New("EventStatus is not a string somehow")
	}

	switch typeValue {
	case Locked.String():
		*t = Locked
	case Unlocked.String():
		*t = Unlocked
	default:
		panic("unknow EventStatus")
	}

	return nil
}

//EventTypes
const (
	Created EventType = iota
	Updated
	Removed

	Locked EventStatus = iota
	Unlocked
)

//PurchaseEvent is an event of purchasing something
type PurchaseEvent struct {
	ID     uint64      `db:"id"`
	Type   EventType   `db:"type"`
	Status EventStatus `db:"status"`
	Entity *Purchase   `db:"payload"`
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

func PurchaseToPB(p *Purchase) pb.Purchase {
	return pb.Purchase{
		Id:       p.ID,
		TotalSum: p.TotalSum,
	}
}

func PBToPurchase(p *pb.Purchase) Purchase {
	return Purchase{
		ID:       p.Id,
		TotalSum: p.TotalSum,
	}
}

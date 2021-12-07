package model

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	pb "github.com/ozonmp/cnm-purchase-api/pkg/cnm-purchase-api"
	"google.golang.org/protobuf/encoding/protojson"
)

// Purchase - purchase entity.
type Purchase struct {
	ID       uint64 `db:"id"`
	TotalSum uint64 `db:"total_sum"`
}

//EventType is type of Event
type EventType uint8

//String converts to string
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

//Scan is a function for custom sqlx diserealisation
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

//String converts to string
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

//Scan is a function for custom sqlx diserealisation
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

//PurchaseHolder holds Purchase. We need this to make custom scan from jsonb to Purchase in PurchaseEvent.
type PurchaseHolder struct {
	Value *Purchase
}

//Scan is a function for custom sqlx diserealisation
func (t *PurchaseHolder) Scan(v interface{}) error {
	log.Error().Msg(fmt.Sprintf("scan purchaseholder"))
	typeValue, ok := v.([]byte)
	if !ok {
		return errors.New("Scan error. interface{} is not a []byte")
	}

	var PBPurchae pb.Purchase

	log.Error().Msg(fmt.Sprintf("unmarhsla protojson"))
	err := protojson.Unmarshal(typeValue, &PBPurchae)
	if err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	log.Error().Msg(fmt.Sprintf("set value"))
	p := PBToPurchase(&PBPurchae)
	t.Value = &p

	log.Error().Msg(fmt.Sprintf("seted value [%+v]", *t.Value))

	return nil
}

//PurchaseEvent is an event of purchasing something
type PurchaseEvent struct {
	ID       uint64         `db:"id"`
	Type     EventType      `db:"type"`
	Status   EventStatus    `db:"status"`
	Purchase PurchaseHolder `db:"payload"`
}

//NewPurchaseEvent creates new PurchaseEvent
func NewPurchaseEvent(
	ID uint64,
	Entity *Purchase,
) *PurchaseEvent {
	return &PurchaseEvent{
		ID:       ID,
		Purchase: PurchaseHolder{Value: Entity},
	}
}

//PurchaseToPB converts Purchase to proto Purchase
func PurchaseToPB(p *Purchase) pb.Purchase {
	return pb.Purchase{
		Id:       p.ID,
		TotalSum: p.TotalSum,
	}
}

//PBToPurchase converts proto Purchase to Purchase
func PBToPurchase(p *pb.Purchase) Purchase {
	return Purchase{
		ID:       p.Id,
		TotalSum: p.TotalSum,
	}
}

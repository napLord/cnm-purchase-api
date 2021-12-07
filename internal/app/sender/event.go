package sender

import (
	"github.com/napLord/cnm-purchase-api/internal/model"
)

//EventSender is an interface of somehting that sends events.
type EventSender interface {
	Send(subdomain *model.PurchaseEvent) error
}

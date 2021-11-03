package sender

import (
	"github.com/napLord/cnm-purchase-api/internal/model"
)

type EventSender interface {
	Send(subdomain *model.PurchaseEvent) error
}

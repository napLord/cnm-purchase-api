package producer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/napLord/cnm-purchase-api/internal/app/remove_queue"
	_ "github.com/napLord/cnm-purchase-api/internal/app/remove_queue"
	"github.com/napLord/cnm-purchase-api/internal/app/sender"
	"github.com/napLord/cnm-purchase-api/internal/model"

	"github.com/napLord/cnm-purchase-api/internal/app/unlock_queue"
)

type Producer interface {
	Start()
}

type producer struct {
	ctx context.Context

	n       uint64
	timeout time.Duration

	sender sender.EventSender
	events <-chan model.PurchaseEvent

	wg   *sync.WaitGroup
	done chan bool

	rq *remove_queue.RemoveQueue
	uq *unlock_queue.UnlockQueue
}

// todo for students: add repo
func NewKafkaProducer(
	ctx context.Context,
	n uint64,
	sender sender.EventSender,
	events <-chan model.PurchaseEvent,
	rq *remove_queue.RemoveQueue,
	uq *unlock_queue.UnlockQueue,
	wg *sync.WaitGroup,
) Producer {
	done := make(chan bool)

	return &producer{
		ctx:    ctx,
		n:      n,
		sender: sender,
		events: events,
		wg:     wg,
		done:   done,
		rq:     rq,
		uq:     uq,
	}
}

func (p *producer) Start() {
	for i := uint64(0); i < p.n; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case event := <-p.events:
					if err := p.sender.Send(&event); err != nil {
						err := p.uq.Unlock(&event)
						if err != nil {
							fmt.Printf(
								"producer can't unlock event[%v] why[%v]",
								event.ID,
								err,
							)
						}
					} else {
						p.rq.Remove(&event)
						if err != nil {
							fmt.Printf(
								"producer can't remove event[%v] why[%v]",
								event.ID,
								err,
							)
						}
					}
				case <-p.ctx.Done():
					return
				}
			}
		}()
	}
}

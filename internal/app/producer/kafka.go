package producer

import (
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
	Close()
}

type producer struct {
	n       uint64
	timeout time.Duration

	sender sender.EventSender
	events <-chan model.PurchaseEvent

	wg sync.WaitGroup

	rq *remove_queue.RemoveQueue
	uq *unlock_queue.UnlockQueue
}

// todo for students: add repo
func NewKafkaProducer(
	n uint64,
	sender sender.EventSender,
	events <-chan model.PurchaseEvent,
	rq *remove_queue.RemoveQueue,
	uq *unlock_queue.UnlockQueue,
) Producer {
	return &producer{
		n:      n,
		sender: sender,
		events: events,
		wg:     sync.WaitGroup{},
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
				case event, ok := <-p.events:
					if !ok {
						return
					}
					if err := p.sender.Send(&event); err != nil {
						p.uq.Unlock(&event)
					} else {
						p.rq.Remove(&event)
					}
				}
			}
		}()
	}
}

func (p *producer) Close() {
	fmt.Printf("producer closing\n")
	p.wg.Wait()
	fmt.Printf("producer closed\n")
}

package producer

import (
	"fmt"
	"sync"

	"github.com/napLord/cnm-purchase-api/internal/app/remove_queue"
	"github.com/napLord/cnm-purchase-api/internal/app/sender"
	"github.com/napLord/cnm-purchase-api/internal/model"

	"github.com/napLord/cnm-purchase-api/internal/app/unlock_queue"
)

//Producer sends event to kafka. unlock events on send error or remove them otherwise
type Producer interface {
	Start()
	Close()
}

type producer struct {
	n uint64

	sender sender.EventSender
	events <-chan model.PurchaseEvent

	wg sync.WaitGroup

	rq *remove_queue.RemoveQueue
	uq *unlock_queue.UnlockQueue
}

//NewKafkaProducer creates new  KafkaProducer
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

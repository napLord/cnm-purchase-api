package producer

import (
	"fmt"
	"sync"
	"time"

	_ "github.com/napLord/cnm-purchase-api/internal/app/remove_queue"
	remove_qeueue "github.com/napLord/cnm-purchase-api/internal/app/remove_queue"
	"github.com/napLord/cnm-purchase-api/internal/app/repo"
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

	wg   *sync.WaitGroup
	done chan bool

	rq *remove_qeueue.RemoveQueue
	uq *unlock_queue.UnlockQueue
}

// todo for students: add repo
func NewKafkaProducer(
	n uint64,
	sender sender.EventSender,
	events <-chan model.PurchaseEvent,
	repo repo.EventRepo,
	removersCount uint64,
	unlockersCount uint64,
) Producer {

	wg := &sync.WaitGroup{}
	done := make(chan bool)

	return &producer{
		n:      n,
		sender: sender,
		events: events,
		wg:     wg,
		done:   done,
		rq:     remove_qeueue.NewRemoveQueue(repo, removersCount),
		uq:     unlock_queue.NewUnlockQueue(repo, unlockersCount),
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
				case <-p.done:
					return
				}
			}
		}()
	}
}

func (p *producer) Close() {
	p.rq.Close()
	p.uq.Close()
	close(p.done)
	p.wg.Wait()
}

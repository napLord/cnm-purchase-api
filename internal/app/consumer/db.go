package consumer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/napLord/cnm-purchase-api/internal/app/repo"
	"github.com/napLord/cnm-purchase-api/internal/model"
)

//Consumer consumes events from db
type Consumer interface {
	Start()
	Close()
}

type consumer struct {
	n      uint64
	events chan<- model.PurchaseEvent

	repo repo.EventRepo

	batchSize uint64
	timeout   time.Duration

	done chan bool
	ctx  context.Context

	wg sync.WaitGroup
}

//NewDbConsumer creates new NewDbConsumer
func NewDbConsumer(
	ctx context.Context,
	n uint64,
	batchSize uint64,
	consumeTimeout time.Duration,
	repo repo.EventRepo,
	events chan<- model.PurchaseEvent,
) Consumer {

	done := make(chan bool)

	return &consumer{
		ctx:       ctx,
		n:         n,
		batchSize: batchSize,
		timeout:   consumeTimeout,
		repo:      repo,
		events:    events,
		wg:        sync.WaitGroup{},
		done:      done,
	}
}

func (c *consumer) Start() {
	for i := uint64(0); i < c.n; i++ {
		c.wg.Add(1)

		go func() {
			defer c.wg.Done()
			ticker := time.NewTicker(c.timeout)
			for {
				select {
				case <-ticker.C:
					events, err := c.repo.Lock(c.batchSize)
					if err != nil {
						fmt.Printf("can't Lock events. why[%+v]", err)
						continue
					}

					for _, event := range events {
						c.events <- event
					}

				case <-c.ctx.Done():
					return
				}
			}
		}()
	}
}

func (c *consumer) Close() {
	fmt.Printf("consumer closing\n")
	c.wg.Wait()
	fmt.Printf("consumer closed\n")
}

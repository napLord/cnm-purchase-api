package unlock_queue

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/napLord/cnm-purchase-api/internal/app/repo"
	"github.com/napLord/cnm-purchase-api/internal/model"
)

const (
	MaxParallelEvents = 16
)

//queue of events to unlock in repo. if event unlock failed, retries it
type UnlockQueue struct {
	ctx context.Context

	repo         repo.EventRepo
	pool         *workerpool.WorkerPool
	running      bool
	workersCount uint64

	wg *sync.WaitGroup

	retryEvents        chan *model.PurchaseEvent
	retryEventsMaxSize uint64
	backslash          uint64

	unlockTimeout time.Duration
}

func NewUnlockQueue(
	ctx context.Context,
	repo repo.EventRepo,
	parallelFactor uint64,
	unlockTimeout time.Duration,
	wg *sync.WaitGroup,
) *UnlockQueue {
	ret := &UnlockQueue{
		ctx:                ctx,
		repo:               repo,
		pool:               workerpool.New(int(parallelFactor)),
		running:            false,
		workersCount:       parallelFactor,
		retryEvents:        make(chan *model.PurchaseEvent, MaxParallelEvents),
		retryEventsMaxSize: MaxParallelEvents,
		backslash:          parallelFactor,
		unlockTimeout:      unlockTimeout,
		wg:                 wg,
	}

	ret.run()

	return ret
}

var ErrQueueIsFull = errors.New("queue is full")

func (q *UnlockQueue) Unlock(e *model.PurchaseEvent) error {
	if !q.running {
		panic("UnlockQueue not running but unlock tryed")
	}

	if uint64(len(q.retryEvents))+q.backslash >= q.retryEventsMaxSize {
		return ErrQueueIsFull
	}

	q.retryEvents <- e

	return nil
}

func (q *UnlockQueue) run() {
	q.running = true

	close := func() {
		q.wg.Add(1)
		defer q.wg.Done()

		_ = <-q.ctx.Done()

		q.running = false
		close(q.retryEvents)
		q.pool.StopWait()

		fmt.Printf("unlock queue closed\n")
	}
	go close()

	for i := 0; i < int(q.workersCount); i++ {
		q.wg.Add(1)

		q.pool.Submit(func() {
			defer q.wg.Done()

			eventsToUnlock := []*model.PurchaseEvent{}

			unlockEvents := func() {
				if len(eventsToUnlock) == 0 {
					return
				}

				IDsToUnlock := make([]uint64, 0, len(eventsToUnlock))

				for i := 0; i < len(eventsToUnlock); i++ {
					IDsToUnlock = append(IDsToUnlock, eventsToUnlock[i].ID)
				}

				err := q.repo.Unlock(IDsToUnlock)

				if err != nil {
					fmt.Printf("can't unlock events[%v] in repo. why[%v]. retry in  queue\n", IDsToUnlock, err)

					for i := 0; i < len(eventsToUnlock); i++ {
						q.retryEvents <- eventsToUnlock[i]
					}
				}

				eventsToUnlock = eventsToUnlock[:0]
			}

			ticker := time.NewTicker(q.unlockTimeout)

			for {

				select {
				case event := <-q.retryEvents:
					if event != nil {
						eventsToUnlock = append(eventsToUnlock, event)
					}

				case <-ticker.C:
					unlockEvents()

				case _ = <-q.ctx.Done():
					unlockEvents()
					return
				}
			}
		})
	}
}

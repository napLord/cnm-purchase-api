package unlock_queue

import (
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
	repo         repo.EventRepo
	pool         *workerpool.WorkerPool
	running      bool
	workersCount uint64

	done chan struct{}
	wg   *sync.WaitGroup

	retryEvents        chan *model.PurchaseEvent
	retryEventsMaxSize uint64
	backslash          uint64

	unlockTimeout time.Duration
}

func NewUnlockQueue(repo repo.EventRepo, parallelFactor uint64, unlockTimeout time.Duration) *UnlockQueue {
	ret := &UnlockQueue{
		repo:               repo,
		pool:               workerpool.New(int(parallelFactor)),
		running:            false,
		workersCount:       parallelFactor,
		done:               make(chan struct{}),
		retryEvents:        make(chan *model.PurchaseEvent, MaxParallelEvents),
		retryEventsMaxSize: MaxParallelEvents,
		backslash:          parallelFactor,
		unlockTimeout:      unlockTimeout,
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

	for i := 0; i < int(q.workersCount); i++ {
		q.pool.Submit(func() {
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

				case _ = <-q.done:
					unlockEvents()
					return
				}
			}
		})
	}
}

func (q *UnlockQueue) Close() {
	q.running = false
	close(q.retryEvents)
	close(q.done)
	q.pool.StopWait()
}

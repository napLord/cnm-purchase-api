package remove_queue

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

//queue of events to remove in repo. if event remove failed, retries it
type RemoveQueue struct {
	repo         repo.EventRepo
	pool         *workerpool.WorkerPool
	running      bool
	workersCount uint64

	done chan struct{}
	wg   *sync.WaitGroup

	retryEvents        chan *model.PurchaseEvent
	retryEventsMaxSize uint64
	backslash          uint64

	removeTimeout time.Duration
}

func NewRemoveQueue(repo repo.EventRepo, parallelFactor uint64, removeTimeout time.Duration) *RemoveQueue {
	ret := &RemoveQueue{
		repo:               repo,
		pool:               workerpool.New(int(parallelFactor)),
		running:            false,
		workersCount:       parallelFactor,
		done:               make(chan struct{}),
		retryEvents:        make(chan *model.PurchaseEvent, MaxParallelEvents),
		retryEventsMaxSize: MaxParallelEvents,
		backslash:          parallelFactor,
		removeTimeout:      removeTimeout,
	}

	ret.run()

	return ret
}

var ErrQueueIsFull = errors.New("queue is full")

func (q *RemoveQueue) Remove(e *model.PurchaseEvent) error {
	if !q.running {
		panic("RemoveQueue not running but remove tryed")
	}

	if uint64(len(q.retryEvents))+q.backslash >= q.retryEventsMaxSize {
		return ErrQueueIsFull
	}

	q.retryEvents <- e

	return nil
}

func (q *RemoveQueue) run() {
	q.running = true

	for i := 0; i < int(q.workersCount); i++ {
		q.pool.Submit(func() {
			eventsToRemove := []*model.PurchaseEvent{}

			removeEvents := func() {
				if len(eventsToRemove) == 0 {
					return
				}

				IDsToRemove := make([]uint64, 0, len(eventsToRemove))

				for i := 0; i < len(eventsToRemove); i++ {
					IDsToRemove = append(IDsToRemove, eventsToRemove[i].ID)
				}

				err := q.repo.Remove(IDsToRemove)

				if err != nil {
					fmt.Printf("can't remove events[%v] in repo. why[%v]. retry in  queue\n", IDsToRemove, err)

					for i := 0; i < len(eventsToRemove); i++ {
						q.retryEvents <- eventsToRemove[i]
					}
				}

				eventsToRemove = eventsToRemove[:0]
			}

			ticker := time.NewTicker(q.removeTimeout)

			for {

				select {
				case event := <-q.retryEvents:
					if event != nil {
						eventsToRemove = append(eventsToRemove, event)
					}

				case <-ticker.C:
					removeEvents()

				case _ = <-q.done:
					removeEvents()
					return
				}
			}
		})
	}
}

func (q *RemoveQueue) Close() {
	q.running = false
	close(q.retryEvents)
	close(q.done)
	q.pool.StopWait()
}

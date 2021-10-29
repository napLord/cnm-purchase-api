package remove_queue

import (
	"errors"
	"fmt"
	"sync"

	"github.com/gammazero/workerpool"
	"github.com/napLord/cnm-purchase-api/internal/app/repo"
	"github.com/napLord/cnm-purchase-api/internal/model"
)

const (
	MaxParallelEvents = 8
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
}

func NewRemoveQueue(repo repo.EventRepo, parallelFactor uint64) *RemoveQueue {
	ret := &RemoveQueue{
		repo:               repo,
		pool:               workerpool.New(int(parallelFactor)),
		running:            false,
		workersCount:       parallelFactor,
		done:               make(chan struct{}),
		retryEvents:        make(chan *model.PurchaseEvent, MaxParallelEvents),
		retryEventsMaxSize: MaxParallelEvents,
		backslash:          parallelFactor,
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
	//add ticker
	q.running = true

	for i := 0; i < int(q.workersCount); i++ {
		q.pool.Submit(func() {
			for {
				select {
				case event := <-q.retryEvents:
					if event != nil {
						err := q.repo.Remove([]uint64{event.ID})

						if err != nil {
							fmt.Printf("can't remove task[%d] in repo. why[%v]. retry in  queue\n", event.ID, err)
							q.retryEvents <- event
						}
					}

				case _ = <-q.done:
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

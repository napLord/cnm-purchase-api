package unlock_queue

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
}

func NewUnlockQueue(repo repo.EventRepo, parallelFactor uint64) *UnlockQueue {
	ret := &UnlockQueue{
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
	//add ticker
	q.running = true

	for i := 0; i < int(q.workersCount); i++ {
		q.pool.Submit(func() {
			for {
				select {
				case event := <-q.retryEvents:
					if event != nil {
						err := q.repo.Unlock([]uint64{event.ID})

						if err != nil {
							fmt.Printf("can't unlock task[%d] in repo. why[%v]. retry in  queue\n", event.ID, err)
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

func (q *UnlockQueue) Close() {
	q.running = false
	close(q.retryEvents)
	close(q.done)
	q.pool.StopWait()
}

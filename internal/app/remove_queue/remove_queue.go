package remove_queue

import (
	"context"
	"fmt"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/napLord/cnm-purchase-api/internal/app/repo"
	"github.com/napLord/cnm-purchase-api/internal/model"
	"go.uber.org/atomic"
)

const (
	//MaxParallelEvents is  how much not unlocked events can Queue hold
	MaxParallelEvents = 16
)

//RemoveQueue is a queue of events to remove in repo. if event remove failed, retries it
type RemoveQueue struct {
	ctx context.Context

	repo         repo.EventRepo
	pool         *workerpool.WorkerPool
	running      *atomic.Bool
	workersCount uint64

	retryEvents        chan *model.PurchaseEvent
	retryEventsMaxSize uint64
	backslash          uint64

	removeTimeout time.Duration
}

//NewRemoveQueue creates new RemoveQueue
func NewRemoveQueue(
	ctx context.Context,
	repo repo.EventRepo,
	parallelFactor uint64,
	removeTimeout time.Duration,
) *RemoveQueue {
	ret := &RemoveQueue{
		ctx:                ctx,
		repo:               repo,
		pool:               workerpool.New(int(parallelFactor)),
		running:            atomic.NewBool(false),
		workersCount:       parallelFactor,
		retryEvents:        make(chan *model.PurchaseEvent, MaxParallelEvents),
		retryEventsMaxSize: MaxParallelEvents,
		backslash:          parallelFactor,
		removeTimeout:      removeTimeout,
	}

	ret.run()

	return ret
}

//Remove event
func (q *RemoveQueue) Remove(e *model.PurchaseEvent) {
	if !q.running.Load() {
		panic("RemoveQueue not running but remove tryed")
	}

	q.retryEvents <- e
}

func (q *RemoveQueue) removeEvents(e []*model.PurchaseEvent) {
	if len(e) == 0 {
		return
	}

	for {
		IDsToRemove := make([]uint64, 0, len(e))

		for i := 0; i < len(e); i++ {
			IDsToRemove = append(IDsToRemove, e[i].ID)
		}

		err := q.repo.Remove(IDsToRemove)

		if err != nil {
			fmt.Printf("can't remove events[%+v] in repo. why[%v]. retry in  queue\n", IDsToRemove, err)

			time.Sleep(q.removeTimeout)
			continue
		}

		return
	}
}

func (q *RemoveQueue) run() {
	q.running.Store(true)

	for i := 0; i < int(q.workersCount); i++ {
		q.pool.Submit(func() {
			eventsToRemove := []*model.PurchaseEvent{}

			ticker := time.NewTicker(q.removeTimeout)

			for {
				select {
				case event, ok := <-q.retryEvents:
					if !ok {
						q.removeEvents(eventsToRemove)

						return
					}

					if event != nil {
						eventsToRemove = append(eventsToRemove, event)
					}

				case <-ticker.C:
					q.removeEvents(eventsToRemove)
					eventsToRemove = eventsToRemove[:0]
				}
			}
		})
	}
}

//Close RemoveQueue
func (q *RemoveQueue) Close() {
	fmt.Printf("remove_queue closing\n")

	_ = <-q.ctx.Done()
	q.running.Store(false)
	close(q.retryEvents)
	q.pool.StopWait()

	fmt.Printf("remove_queue closed\n")
}

package unlock_queue

import (
	"context"
	"fmt"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/ozonmp/cnm-purchase-api/internal/app/repo"
	"github.com/ozonmp/cnm-purchase-api/internal/model"
	"go.uber.org/atomic"
)

const (
	//MaxParallelEvents is  how much not unlocked events can Queue hold
	MaxParallelEvents = 16
)

//UnlockQueue is a queue of events to unlock in repo. if event unlock failed, retries it
type UnlockQueue struct {
	ctx context.Context

	repo         repo.EventRepo
	pool         *workerpool.WorkerPool
	running      *atomic.Bool
	workersCount uint64

	retryEvents        chan *model.PurchaseEvent
	retryEventsMaxSize uint64
	backslash          uint64

	unlockTimeout time.Duration
}

//NewUnlockQueue creates new UnlockQueue
func NewUnlockQueue(
	ctx context.Context,
	repo repo.EventRepo,
	parallelFactor uint64,
	unlockTimeout time.Duration,
) *UnlockQueue {
	ret := &UnlockQueue{
		ctx:                ctx,
		repo:               repo,
		pool:               workerpool.New(int(parallelFactor)),
		running:            atomic.NewBool(false),
		workersCount:       parallelFactor,
		retryEvents:        make(chan *model.PurchaseEvent, MaxParallelEvents),
		retryEventsMaxSize: MaxParallelEvents,
		backslash:          parallelFactor,
		unlockTimeout:      unlockTimeout,
	}

	ret.run()

	return ret
}

//Unlock event
func (q *UnlockQueue) Unlock(e *model.PurchaseEvent) {
	if !q.running.Load() {
		panic("UnlockQueue not running but unlock tryed")
	}

	q.retryEvents <- e
}

func (q *UnlockQueue) unlockEvents(e []*model.PurchaseEvent) {
	if len(e) == 0 {
		return
	}

	for {
		IDsToUnlock := make([]uint64, 0, len(e))

		for i := 0; i < len(e); i++ {
			IDsToUnlock = append(IDsToUnlock, e[i].ID)
		}

		err := q.repo.Unlock(IDsToUnlock)

		if err != nil {
			fmt.Printf("can't unlock events[%+v] in repo. why[%v]. retry in  queue\n", IDsToUnlock, err)

			time.Sleep(q.unlockTimeout)
			continue
		}

		return
	}
}

func (q *UnlockQueue) run() {
	q.running.Store(true)

	for i := 0; i < int(q.workersCount); i++ {
		q.pool.Submit(func() {
			eventsToUnlock := []*model.PurchaseEvent{}

			ticker := time.NewTicker(q.unlockTimeout)

			for {
				select {
				case event, ok := <-q.retryEvents:
					if !ok {
						q.unlockEvents(eventsToUnlock)

						return
					}

					if event != nil {
						eventsToUnlock = append(eventsToUnlock, event)
					}

				case <-ticker.C:
					q.unlockEvents(eventsToUnlock)
					eventsToUnlock = eventsToUnlock[:0]
				}
			}
		})
	}
}

//Close UnlockQueue
func (q *UnlockQueue) Close() {
	fmt.Printf("unlock_queue closing\n")

	_ = <-q.ctx.Done()
	q.running.Store(false)
	close(q.retryEvents)
	q.pool.StopWait()

	fmt.Printf("unlock_queue closed\n")
}

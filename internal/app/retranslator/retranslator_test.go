package retranslator

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/napLord/cnm-purchase-api/internal/app/remove_queue"
	"github.com/napLord/cnm-purchase-api/internal/app/unlock_queue"
	"github.com/napLord/cnm-purchase-api/internal/mocks"
	"github.com/napLord/cnm-purchase-api/internal/model"
	"github.com/stretchr/testify/assert"
)

var testTimeout = time.Second * 6

func TestStart(t *testing.T) {
	//pre
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockEventRepo(ctrl)
	sender := mocks.NewMockEventSender(ctrl)

	cfg := Config{
		ChannelSize:    512,
		ConsumerCount:  2,
		ConsumeSize:    10,
		ConsumeTimeout: 400 * time.Millisecond,
		ProducerCount:  2,
		WorkerCount:    1,
		Repo:           repo,
		Sender:         sender,
		removeTimeout:  time.Second,
		unlockTimeout:  time.Second,
	}
	eventsCount := 16

	purchase1 := model.Purchase{1, 1}

	goodSendEvents := map[uint64]struct{}{}
	badSendEvents := map[uint64]struct{}{}

	removedEvents := map[uint64]struct{}{}
	unlockedEvents := map[uint64]struct{}{}

	events := []model.PurchaseEvent{}
	eventsLockIdx := uint64(0)

	for i := 0; i < eventsCount; i++ {
		if i%2 == 0 {
			badSendEvents[uint64(i)] = struct{}{}
		} else {
			goodSendEvents[uint64(i)] = struct{}{}
		}

		events = append(events, *model.NewPurchaseEvent(uint64(i), &purchase1))
	}

	//expect behaviour

	//expects lock. returns eventsCount events
	repo.EXPECT().
		Lock(gomock.Any()).
		Times(eventsCount).
		DoAndReturn(
			func(n uint64) ([]model.PurchaseEvent, error) {
				ret := []model.PurchaseEvent{events[eventsLockIdx]}

				atomic.AddUint64(&eventsLockIdx, 1)

				return ret, nil
			},
		)

	//expects locks after
	repo.EXPECT().Lock(gomock.Any()).AnyTimes()

	//expects eventsCount sends.
	//return err if events is bad or nil if event is good
	sendBad := sender.EXPECT().
		Send(gomock.Any()).
		AnyTimes().
		DoAndReturn(func(e *model.PurchaseEvent) error {
			fmt.Printf("called Send! with e[%v]\n", e)

			if _, ok := goodSendEvents[e.ID]; !ok {
				return errors.New("send failed. bad event")
			}

			return nil
		})

	sender.EXPECT().Send(gomock.Any()).AnyTimes().After(sendBad)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	//expects removes as we got good sends
	removeCalls := repo.EXPECT().
		Remove(gomock.Any()).
		AnyTimes().
		Do(func(eventIDs []uint64) error {
			fmt.Printf("called Remove!with ids[%v]\n", eventIDs)

			for _, k := range eventIDs {
				removedEvents[k] = struct{}{}
			}

			if len(removedEvents) == len(goodSendEvents) {
				assert.Equal(t, removedEvents, goodSendEvents)
				defer wg.Done()
			}

			return nil
		})

	repo.EXPECT().Remove(gomock.Any()).AnyTimes().After(removeCalls)

	//expects  unlocks as we got bad sends
	unlockCalls := repo.EXPECT().
		Unlock(gomock.Any()).
		AnyTimes().
		Do(func(eventIDs []uint64) error {
			fmt.Printf("called Unlock!with ids[%v]\n", eventIDs)

			for _, k := range eventIDs {
				unlockedEvents[k] = struct{}{}
			}

			if len(unlockedEvents) == len(badSendEvents) {
				assert.Equal(t, unlockedEvents, badSendEvents)
				defer wg.Done()
			}

			return nil
		})

	repo.EXPECT().Unlock(gomock.Any()).AnyTimes().After(unlockCalls)

	retranslator := NewRetranslator(cfg)
	retranslator.Start()
	defer retranslator.Close()

	wg.Wait()
}

func TestBrokenDBUnlock(t *testing.T) {
	//pre
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockEventRepo(ctrl)
	sender := mocks.NewMockEventSender(ctrl)

	cfg := Config{
		ChannelSize:    512,
		ConsumerCount:  2,
		ConsumeSize:    10,
		ConsumeTimeout: 50 * time.Millisecond,
		ProducerCount:  2,
		WorkerCount:    2,
		Repo:           repo,
		Sender:         sender,
		removeTimeout:  200 * time.Millisecond,
		unlockTimeout:  200 * time.Millisecond,
	}

	purchase1 := model.Purchase{1, 1}

	eventsCount := unlock_queue.MaxParallelEvents / 2

	unlockFailedCount := 10

	unlockedEvents := map[uint64]struct{}{}

	events := []model.PurchaseEvent{}
	eventsID := map[uint64]struct{}{}
	eventsLockIdx := uint64(0)

	for i := 0; i < eventsCount; i++ {
		events = append(events, *model.NewPurchaseEvent(uint64(i), &purchase1))

		eventsID[uint64(i)] = struct{}{}
	}

	//expect behaviour

	//lock 1 event eventsCount times. we expect them to be unlocked further due to bad sends
	lockWithGoodEv := repo.EXPECT().
		Lock(gomock.Any()).
		Times(eventsCount).
		DoAndReturn(
			func(n uint64) ([]model.PurchaseEvent, error) {
				ret := []model.PurchaseEvent{events[eventsLockIdx]}

				atomic.AddUint64(&eventsLockIdx, 1)

				return ret, nil
			},
		)

	//expects locks after
	repo.EXPECT().Lock(gomock.Any()).AnyTimes().After(lockWithGoodEv)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	//expects eventsCount sends. sends are bad so return with error
	sendBad := sender.EXPECT().
		Send(gomock.Any()).
		Times(eventsCount).
		DoAndReturn(func(e *model.PurchaseEvent) error {
			fmt.Printf("called Send! with e[%v]. responding with err\n", e)

			return errors.New("send failed")
		})

	sender.EXPECT().Send(gomock.Any()).AnyTimes().After(sendBad)

	//expects Unlocks of events.
	//emulating broken unlock. so respons with error unlockFailedCount times
	unlockCallsBad := repo.EXPECT().
		Unlock(gomock.Any()).
		Times(unlockFailedCount).
		DoAndReturn(func(eventIDs []uint64) error {
			fmt.Printf("called Unlock!with ids[%v]. responding with error\n", eventIDs)

			return errors.New(fmt.Sprintf("db can't remove events id[%v]\n", eventIDs))
		})

	//extepcts good eventsCount unlocks
	unlockCallsGood := repo.EXPECT().
		Unlock(gomock.Any()).
		AnyTimes().
		After(unlockCallsBad).
		DoAndReturn(func(eventIDs []uint64) error {
			fmt.Printf("called Unlock!with ids[%v].\n", eventIDs)

			for _, k := range eventIDs {
				unlockedEvents[k] = struct{}{}
			}

			fmt.Printf("unlockedEvents[%v]\n", eventIDs)

			if len(unlockedEvents) == len(eventsID) {
				assert.Equal(t, unlockedEvents, eventsID)
				defer wg.Done()
			}

			return nil
		})

	repo.EXPECT().Unlock(gomock.Any()).AnyTimes().After(unlockCallsGood).After(unlockCallsBad)

	retranslator := NewRetranslator(cfg)
	retranslator.Start()
	defer retranslator.Close()

	wg.Wait()
}

//copypaste from TestBrokenDBUnlock
func TestBrokenDBRemove(t *testing.T) {
	//pre
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockEventRepo(ctrl)
	sender := mocks.NewMockEventSender(ctrl)

	cfg := Config{
		ChannelSize:    512,
		ConsumerCount:  2,
		ConsumeSize:    10,
		ConsumeTimeout: 50 * time.Millisecond,
		ProducerCount:  2,
		WorkerCount:    2,
		Repo:           repo,
		Sender:         sender,
		removeTimeout:  200 * time.Millisecond,
		unlockTimeout:  200 * time.Millisecond,
	}

	purchase1 := model.Purchase{1, 1}

	eventsCount := remove_queue.MaxParallelEvents / 2

	removeFailedCount := 10

	removeedEvents := map[uint64]struct{}{}

	events := []model.PurchaseEvent{}
	eventsID := map[uint64]struct{}{}
	eventsLockIdx := uint64(0)

	for i := 0; i < eventsCount; i++ {
		events = append(events, *model.NewPurchaseEvent(uint64(i), &purchase1))

		eventsID[uint64(i)] = struct{}{}
	}

	//expect behaviour

	//lock 1 event eventsCount times. we expect them to be further removeed due to bad sends
	lockWithGoodEv := repo.EXPECT().
		Lock(gomock.Any()).
		Times(eventsCount).
		DoAndReturn(
			func(n uint64) ([]model.PurchaseEvent, error) {
				ret := []model.PurchaseEvent{events[eventsLockIdx]}

				atomic.AddUint64(&eventsLockIdx, 1)

				return ret, nil
			},
		)

	//expects locks after
	repo.EXPECT().Lock(gomock.Any()).AnyTimes().After(lockWithGoodEv)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	//expects eventsCount sends. sends are bad so return with error
	sendBad := sender.EXPECT().
		Send(gomock.Any()).
		Times(eventsCount).
		DoAndReturn(func(e *model.PurchaseEvent) error {
			fmt.Printf("called Send! with e[%v]. responding good\n", e)

			return nil
		})

	sender.EXPECT().Send(gomock.Any()).AnyTimes().After(sendBad)

	//expects Removes of events.
	//emulating broken remove. so respons with error removeFailedCount times
	removeCallsBad := repo.EXPECT().
		Remove(gomock.Any()).
		Times(removeFailedCount).
		DoAndReturn(func(eventIDs []uint64) error {
			fmt.Printf("called Remove!with ids[%v]. responding with error\n", eventIDs)

			return errors.New(fmt.Sprintf("db can't remove events id[%v]\n", eventIDs))
		})

	//extepcts good eventsCount removes
	removeCallsGood := repo.EXPECT().
		Remove(gomock.Any()).
		AnyTimes().
		After(removeCallsBad).
		DoAndReturn(func(eventIDs []uint64) error {
			fmt.Printf("called Remove!with ids[%v].\n", eventIDs)

			for _, k := range eventIDs {
				removeedEvents[k] = struct{}{}
			}

			fmt.Printf("removeedEvents[%v]\n", eventIDs)

			if len(removeedEvents) == len(eventsID) {
				assert.Equal(t, removeedEvents, eventsID)
				defer wg.Done()
			}

			return nil
		})

	repo.EXPECT().Remove(gomock.Any()).AnyTimes().After(removeCallsGood).After(removeCallsBad)

	retranslator := NewRetranslator(cfg)
	retranslator.Start()
	defer retranslator.Close()

	wg.Wait()
}

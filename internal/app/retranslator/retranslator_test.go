package retranslator

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/napLord/cnm-purchase-api/internal/app/remove_queue"
	"github.com/napLord/cnm-purchase-api/internal/app/unlock_queue"
	"github.com/napLord/cnm-purchase-api/internal/mocks"
	"github.com/napLord/cnm-purchase-api/internal/model"
)

func TestStart(t *testing.T) {
	//test timeout
	go func() {
		time.Sleep(time.Second * 10)

		panic("test timeout")
	}()

	//pre
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockEventRepo(ctrl)
	sender := mocks.NewMockEventSender(ctrl)

	cfg := Config{
		ChannelSize:    512,
		ConsumerCount:  2,
		ConsumeSize:    10,
		ConsumeTimeout: 100 * time.Millisecond,
		ProducerCount:  2,
		WorkerCount:    2,
		Repo:           repo,
		Sender:         sender,
	}
	eventsCount := 10

	purchase1 := model.Purchase{1, 1}
	purchase2 := model.Purchase{2, 1}

	eventGoodID := uint64(1)
	eventBadID := uint64(2)

	eventGoodSend := model.NewPurchaseEvent(eventGoodID, &purchase1)
	eventBadSend := model.NewPurchaseEvent(eventBadID, &purchase2)

	//expect behaviour

	//expects eventsCount/2 locks. returns events that won't fail on send
	repo.EXPECT().
		Lock(gomock.Any()).
		Times(eventsCount / 2).
		DoAndReturn(
			func(n uint64) ([]model.PurchaseEvent, error) {
				return []model.PurchaseEvent{*eventGoodSend}, nil
			},
		)

	//expects eventsCount/2 locks. returns events that WILL fail on send
	repo.EXPECT().
		Lock(gomock.Any()).
		Times(eventsCount / 2).
		DoAndReturn(
			func(n uint64) ([]model.PurchaseEvent, error) {
				return []model.PurchaseEvent{*eventBadSend}, nil
			},
		)

	//expects locks after
	repo.EXPECT().Lock(gomock.Any()).AnyTimes()

	wg := &sync.WaitGroup{}
	wg.Add(eventsCount) //generates

	//expects eventsCount/2 sends with bad events. return with error
	sendBad := sender.EXPECT().
		Send(eventBadSend).
		Times(eventsCount / 2).
		DoAndReturn(func(e *model.PurchaseEvent) error {
			fmt.Printf("called Send! with e[%v]. responding with err\n", e)

			return errors.New("send failed")
		})

	//expects eventsCount/2 sends with good events. return nil
	sendGood := sender.EXPECT().
		Send(eventGoodSend).
		Times(eventsCount / 2).
		DoAndReturn(func(e *model.PurchaseEvent) error {
			fmt.Printf("called Send! with e[%v]. responding good\n", e)

			return nil
		})

	sender.EXPECT().Send(gomock.Any()).AnyTimes().After(sendBad).After(sendGood)

	//expects eventsCount/2 removes as we got eventsCount/2 good sends
	removeCalls := repo.EXPECT().
		Remove([]uint64{eventGoodID}).
		Times(eventsCount / 2).
		Do(func(eventIDs []uint64) error {
			wg.Done()

			fmt.Printf("called Remove!with ids[%v]\n", eventIDs)

			return nil
		})

	repo.EXPECT().Remove(gomock.Any()).AnyTimes().After(removeCalls)

	//expects eventsCount/2 unlocks as we got eventsCount/2 bad sends
	unlockCalls := repo.EXPECT().
		Unlock([]uint64{eventBadID}).
		Times(eventsCount / 2).
		Do(func(eventIDs []uint64) error {
			wg.Done()

			fmt.Printf("called Unlock!with ids[%v]\n", eventIDs)

			return nil
		})

	repo.EXPECT().Unlock(gomock.Any()).AnyTimes().After(unlockCalls)

	retranslator := NewRetranslator(cfg)
	retranslator.Start()
	defer retranslator.Close()

	wg.Wait()
}

func TestBrokenDBUnlock(t *testing.T) {
	//test timeout
	go func() {
		time.Sleep(time.Second * 5)

		panic("test timeout")
	}()

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
	}

	purchase1 := model.Purchase{1, 1}
	eventID := uint64(1)
	event1 := model.NewPurchaseEvent(eventID, &purchase1)

	eventsCount := unlock_queue.MaxParallelEvents / 2

	//expect behaviour

	//lock 1 event eventsCount times. we expect them to be unlocked due to bad sends
	lockWithGoodEv := repo.EXPECT().
		Lock(gomock.Any()).
		Times(eventsCount).
		DoAndReturn(
			func(n uint64) ([]model.PurchaseEvent, error) {
				return []model.PurchaseEvent{*event1}, nil
			},
		)

	//expects locks after
	repo.EXPECT().Lock(gomock.Any()).AnyTimes().After(lockWithGoodEv)

	wg := &sync.WaitGroup{}
	wg.Add(eventsCount)

	//expects eventsCount sends. sends are bad so return with error
	sendBad := sender.EXPECT().
		Send(event1).
		Times(eventsCount).
		DoAndReturn(func(e *model.PurchaseEvent) error {
			fmt.Printf("called Send! with e[%v]. responding with err\n", e)

			return errors.New("send failed")
		})

	sender.EXPECT().Send(gomock.Any()).AnyTimes().After(sendBad)

	//expects Unlocks of events.
	//emulating broken unlock. so respons with error eventsCout * 2 times
	unlockCallsBad := repo.EXPECT().
		Unlock([]uint64{eventID}).
		Times(eventsCount * 2).
		DoAndReturn(func(eventIDs []uint64) error {
			fmt.Printf("called Unlock!with ids[%v]. responding with error\n", eventIDs)

			return errors.New(fmt.Sprintf("db can't remove events id[%v]\n", eventIDs))
		})

	//extepcts good eventsCount unlocks
	unlockCallsGood := repo.EXPECT().
		Unlock([]uint64{eventID}).
		Times(eventsCount).
		After(unlockCallsBad).
		DoAndReturn(func(eventIDs []uint64) error {
			wg.Done()

			fmt.Printf("called Unlock!with ids[%v]\n", eventIDs)

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
	//test timeout
	go func() {
		time.Sleep(time.Second * 5)

		panic("test timeout")
	}()

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
	}

	purchase1 := model.Purchase{1, 1}
	eventID := uint64(1)
	event1 := model.NewPurchaseEvent(eventID, &purchase1)

	eventsCount := remove_queue.MaxParallelEvents / 2

	//expect behaviour

	//lock 1 event eventsCount times. we expect them to be removed due to good sends
	lockWithGoodEv := repo.EXPECT().
		Lock(gomock.Any()).
		Times(eventsCount).
		DoAndReturn(
			func(n uint64) ([]model.PurchaseEvent, error) {
				return []model.PurchaseEvent{*event1}, nil
			},
		)

	//expects locks after
	repo.EXPECT().Lock(gomock.Any()).AnyTimes().After(lockWithGoodEv)

	wg := &sync.WaitGroup{}
	wg.Add(eventsCount)

	//expects eventsCount sends. sends are bad so return with error
	sendBad := sender.EXPECT().
		Send(event1).
		Times(eventsCount).
		DoAndReturn(func(e *model.PurchaseEvent) error {
			fmt.Printf("called Send! with e[%v]\n", e)

			return nil
		})

	sender.EXPECT().Send(gomock.Any()).AnyTimes().After(sendBad)

	//expects Removes of events.
	//emulating broken Remove. so respons with error eventsCount * 2 times
	removeCallsBad := repo.EXPECT().
		Remove([]uint64{eventID}).
		Times(eventsCount * 2).
		DoAndReturn(func(eventIDs []uint64) error {
			fmt.Printf("called Remove! with ids[%v]. responding with error\n", eventIDs)

			return errors.New(fmt.Sprintf("db can't remove events id[%v]\n", eventIDs))
		})

	//expects good eventsCount removes
	removeCallsGood := repo.EXPECT().
		Remove([]uint64{eventID}).
		Times(eventsCount).
		After(removeCallsBad).
		DoAndReturn(func(eventIDs []uint64) error {
			wg.Done()

			fmt.Printf("called Remove!with ids[%v]\n", eventIDs)

			return nil
		})

	repo.EXPECT().Remove(gomock.Any()).AnyTimes().After(removeCallsGood).After(removeCallsBad)

	retranslator := NewRetranslator(cfg)
	retranslator.Start()
	defer retranslator.Close()

	wg.Wait()
}

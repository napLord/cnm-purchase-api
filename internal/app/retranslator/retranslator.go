package retranslator

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/napLord/cnm-purchase-api/internal/app/consumer"
	"github.com/napLord/cnm-purchase-api/internal/app/producer"
	"github.com/napLord/cnm-purchase-api/internal/app/remove_queue"
	"github.com/napLord/cnm-purchase-api/internal/app/repo"
	"github.com/napLord/cnm-purchase-api/internal/app/sender"
	"github.com/napLord/cnm-purchase-api/internal/app/unlock_queue"
	"github.com/napLord/cnm-purchase-api/internal/model"
)

type Retranslator interface {
	Start()
	Close()
}

type Config struct {
	ChannelSize uint64

	ConsumerCount  uint64
	ConsumeSize    uint64
	ConsumeTimeout time.Duration

	ProducerCount uint64
	WorkerCount   int

	Repo   repo.EventRepo
	Sender sender.EventSender

	removeTimeout time.Duration
	unlockTimeout time.Duration
}

type retranslator struct {
	events   chan model.PurchaseEvent
	consumer consumer.Consumer
	producer producer.Producer
	rq       *remove_queue.RemoveQueue
	uq       *unlock_queue.UnlockQueue

	cancelFunc context.CancelFunc
	wg         *sync.WaitGroup
}

func NewRetranslator(cfg Config) Retranslator {
	events := make(chan model.PurchaseEvent, cfg.ChannelSize)

	wg := &sync.WaitGroup{}
	ctx, cancelFunc := context.WithCancel(context.Background())

	consumer := consumer.NewDbConsumer(
		ctx,
		cfg.ConsumerCount,
		cfg.ConsumeSize,
		cfg.ConsumeTimeout,
		cfg.Repo,
		events,
		wg,
	)

	remove_queue := remove_queue.NewRemoveQueue(
		ctx,
		cfg.Repo,
		uint64(cfg.WorkerCount),
		cfg.removeTimeout,
		wg,
	)

	unlock_queue := unlock_queue.NewUnlockQueue(
		ctx,
		cfg.Repo,
		uint64(cfg.WorkerCount),
		cfg.unlockTimeout,
		wg,
	)

	producer := producer.NewKafkaProducer(
		ctx,
		cfg.ProducerCount,
		cfg.Sender,
		events,
		remove_queue,
		unlock_queue,
		wg,
	)

	return &retranslator{
		events:     events,
		consumer:   consumer,
		producer:   producer,
		rq:         remove_queue,
		uq:         unlock_queue,
		wg:         wg,
		cancelFunc: cancelFunc,
	}
}

func (r *retranslator) Start() {
	fmt.Printf("retranslator starts\n")
	r.producer.Start()
	r.consumer.Start()
}

func (r *retranslator) Close() {
	r.cancelFunc()

	r.wg.Wait()

	fmt.Printf("retranslator closed\n")
}

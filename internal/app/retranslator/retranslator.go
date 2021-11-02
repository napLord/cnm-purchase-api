package retranslator

import (
	"context"
	"fmt"
	"time"

	"github.com/napLord/cnm-purchase-api/internal/app/closer"
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

	ctx context.Context
}

func NewRetranslator(ctx context.Context, clr *closer.Closer, cfg Config) Retranslator {
	events := make(chan model.PurchaseEvent, cfg.ChannelSize)

	consumer := consumer.NewDbConsumer(
		ctx,
		cfg.ConsumerCount,
		cfg.ConsumeSize,
		cfg.ConsumeTimeout,
		cfg.Repo,
		events,
	)

	consumerCh := make(chan struct{})
	producerCh := make(chan struct{})

	clr.Add(func() {
		consumer.Close()
		close(consumerCh)
	})

	remove_queue := remove_queue.NewRemoveQueue(
		ctx,
		cfg.Repo,
		uint64(cfg.WorkerCount),
		cfg.removeTimeout,
	)
	clr.Add(func() {
		<-producerCh
		remove_queue.Close()
	})

	unlock_queue := unlock_queue.NewUnlockQueue(
		ctx,
		cfg.Repo,
		uint64(cfg.WorkerCount),
		cfg.unlockTimeout,
	)
	clr.Add(func() {
		<-producerCh
		unlock_queue.Close()
	})

	producer := producer.NewKafkaProducer(
		cfg.ProducerCount,
		cfg.Sender,
		events,
		remove_queue,
		unlock_queue,
	)
	clr.Add(func() {
		<-consumerCh
		close(events)
		producer.Close()
		close(producerCh)
	})

	return &retranslator{
		events:   events,
		consumer: consumer,
		producer: producer,
		rq:       remove_queue,
		uq:       unlock_queue,
	}
}

func (r *retranslator) Start() {
	fmt.Printf("retranslator starts\n")
	r.producer.Start()
	r.consumer.Start()
}

func (r *retranslator) Close() {
	fmt.Printf("retranslator closed\n")
}

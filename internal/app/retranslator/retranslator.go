package retranslator

import (
	"fmt"
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
}

func NewRetranslator(cfg Config) Retranslator {
	events := make(chan model.PurchaseEvent, cfg.ChannelSize)

	consumer := consumer.NewDbConsumer(
		cfg.ConsumerCount,
		cfg.ConsumeSize,
		cfg.ConsumeTimeout,
		cfg.Repo,
		events)

	remove_queue := remove_queue.NewRemoveQueue(cfg.Repo, uint64(cfg.WorkerCount), cfg.removeTimeout)
	unlock_queue := unlock_queue.NewUnlockQueue(cfg.Repo, uint64(cfg.WorkerCount), cfg.unlockTimeout)

	producer := producer.NewKafkaProducer(
		cfg.ProducerCount,
		cfg.Sender,
		events,
		uint64(cfg.WorkerCount),
		uint64(cfg.WorkerCount),
		remove_queue,
		unlock_queue,
	)

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
	fmt.Printf("retranslator closes\n")
	r.consumer.Close()
	fmt.Printf("consumer closed\n")
	r.rq.Close()
	fmt.Printf("remove queue closed\n")
	r.uq.Close()
	fmt.Printf("unlock queue closed\n")
	r.producer.Close()
	fmt.Printf("producer closed\n")
}

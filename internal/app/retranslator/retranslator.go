package retranslator

import (
	"fmt"
	"time"

	"github.com/napLord/cnm-purchase-api/internal/app/consumer"
	"github.com/napLord/cnm-purchase-api/internal/app/producer"
	"github.com/napLord/cnm-purchase-api/internal/app/repo"
	"github.com/napLord/cnm-purchase-api/internal/app/sender"
	"github.com/napLord/cnm-purchase-api/internal/model"

	"github.com/gammazero/workerpool"
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
}

type retranslator struct {
	events     chan model.PurchaseEvent
	consumer   consumer.Consumer
	producer   producer.Producer
	workerPool *workerpool.WorkerPool
}

func NewRetranslator(cfg Config) Retranslator {
	events := make(chan model.PurchaseEvent, cfg.ChannelSize)
	workerPool := workerpool.New(cfg.WorkerCount)

	consumer := consumer.NewDbConsumer(
		cfg.ConsumerCount,
		cfg.ConsumeSize,
		cfg.ConsumeTimeout,
		cfg.Repo,
		events)

	producer := producer.NewKafkaProducer(
		cfg.ProducerCount,
		cfg.Sender,
		events,
		cfg.Repo,
		uint64(cfg.WorkerCount),
		uint64(cfg.WorkerCount),
	)

	return &retranslator{
		events:     events,
		consumer:   consumer,
		producer:   producer,
		workerPool: workerPool,
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
	r.producer.Close()
	fmt.Printf("producer closed\n")
	r.workerPool.StopWait()
	fmt.Printf("workerpoold closed\n")
}

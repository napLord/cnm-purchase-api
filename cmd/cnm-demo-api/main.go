package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/napLord/cnm-purchase-api/internal/app/closer"       //nolint
	"github.com/napLord/cnm-purchase-api/internal/app/retranslator" //nolint
	"golang.org/x/net/context"
	//"github.com/ozonmp/cnm-demo-api/internal/app/retranslator"
)

func main() {
	sigs := make(chan os.Signal, 1)

	cfg := retranslator.Config{
		ChannelSize:   512,
		ConsumerCount: 2,
		ConsumeSize:   10,
		ProducerCount: 28,
		WorkerCount:   2,
	}

	ctx, cancel := context.WithCancel(context.Background())

	clr := closer.NewCloser()

	clr.Add(func() {
		cancel()
	})
	defer clr.RunAll()

	retranslator := retranslator.NewRetranslator(ctx, clr, cfg)
	retranslator.Start()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
}

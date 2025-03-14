package main

import (
	"MP4Processor/config"
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	channelBufferSize = 1024
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown on SIGINT or SIGTERM
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		cancel()
	}()

	conf, err := config.LoadConfig()
	nc, err := nats.Connect(conf.NATS.URL)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to NATS: %v", err))
	}
	defer nc.Close()

	msgChan := make(chan *nats.Msg, channelBufferSize) // Buffered channel to avoid blocking
	sub, err := nc.ChanSubscribe(conf.NATS.Mp4FilePathsTopic, msgChan)
	if err != nil {
		panic(fmt.Sprintf("failed to subscribe to topic %s: %v", conf.NATS.Mp4FilePathsTopic, err))
	}
	defer drainAndUnsubscribe(sub)

	if err := process(ctx, msgChan, conf.File.OutputPath, conf.NATS.ProcessResultTopic, nc.Publish); err != nil {
		log.Fatalf("Process failed: %v", err)
	}
}

func drainAndUnsubscribe(sub *nats.Subscription) {
	if err := sub.Drain(); err != nil {
		log.Printf("Error draining topic: %v\n", err)
	}
	if err := sub.Unsubscribe(); err != nil {
		log.Printf("Error unsubscribing topic: %v\n", err)
	}
}

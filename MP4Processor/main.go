package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/nats-io/nats.go"

	"MP4Processor/config"
	"MP4Processor/model"
	"MP4Processor/processor"
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
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	nc, err := nats.Connect(conf.NATS.URL)
	if err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	msgChan := make(chan *nats.Msg, conf.NATS.BufferSize) // Buffered channel to avoid blocking
	sub, err := nc.ChanSubscribe(conf.NATS.Mp4FilePathsTopic, msgChan)
	if err != nil {
		log.Fatalf("failed to subscribe to topic %s: %v", conf.NATS.Mp4FilePathsTopic, err)
	}
	defer drainAndUnsubscribe(sub)

	// Setup processor
	p := processor.New(conf.File.OutputPath)

	var wg sync.WaitGroup
	defer wg.Wait()

	// Listen to nats message
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down...")
			return
		case msg, ok := <-msgChan:
			if !ok {
				log.Println("Nats message channel closed")
				return
			}

			inputFile := string(msg.Data)

			wg.Add(1)
			go func(p *processor.Processor, inputFile string) {
				defer wg.Done()
				// Process incoming message
				natsRes := processNatsMessage(p, inputFile, conf.File.OutputPath)

				// Marshal nats response
				res, err := json.Marshal(natsRes)
				if err != nil {
					log.Printf("failed to marshal nats response: %v", err)
					return
				}

				// Publish response
				if err = nc.Publish(conf.NATS.ProcessResultTopic, res); err != nil {
					log.Printf("Error publishing message: %v\n", err)
				}
			}(p, inputFile)
		}
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

func processNatsMessage(p *processor.Processor, inputFile string, outputPath string) (natsRes model.ProcessedFileMessage) {
	natsRes = model.ProcessedFileMessage{
		FileName:   inputFile,
		StatusCode: model.StatusSuccessful,
		Message:    "File processed successfully",
		ResultPath: outputPath,
	}
	var pErr error // Process error
	pErr = p.ProcessMP4(inputFile)
	if pErr != nil {
		natsRes.StatusCode = model.StatusFailed
		natsRes.Message = pErr.Error()
		natsRes.ResultPath = ""
	}

	return natsRes
}

package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const (
	channelBufferSize = 1024
)

type PublishFunc func(subject string, msg []byte) error

func process(ctx context.Context, msgChan chan *nats.Msg, processResultTopic string, publish PublishFunc) error {
	var wg sync.WaitGroup
	defer wg.Wait()

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down...")
			return nil
		case msg, ok := <-msgChan:
			if !ok {
				log.Println("Message channel closed")
				return nil
			}

			wg.Add(1)
			go func(msg *nats.Msg) {
				defer wg.Done()
				processMessage(msg, processResultTopic, publish)
			}(msg)
		}
	}
}

func processMessage(msg *nats.Msg, processResultTopic string, publish PublishFunc) {
	resultMessage := handleExtractionMessage(msg)
	byteArray, err := json.Marshal(resultMessage)
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return
	}
	if err := publish(processResultTopic, byteArray); err != nil {
		fmt.Println("Error publishing message:", err)
	}
}

func handleExtractionMessage(msg *nats.Msg) *processedFileMessage {
	filePath := string(msg.Data)
	fmt.Printf("Received file path: %s\n", filePath)

	boxes, err := ExtractInitializationSegment(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return &processedFileMessage{
			FileName:   filePath,
			StatusCode: Status(StatusFailedProcessing),
			Message:    fmt.Sprintf("Failed to process init segment: %v", err),
			ResultPath: "",
		}
	}

	for _, box := range boxes {
		fmt.Printf("MP4Box Type: %s, Size: %d\n", box.Type, box.Size)
	}

	resultPath, err := writeResultIntoFile("output.mp4", boxes)

	if err != nil {
		return &processedFileMessage{
			FileName:   filePath,
			StatusCode: Status(StatusFailedProcessing),
			Message:    fmt.Sprintf("Failed to write into a file: %v", err),
			ResultPath: "",
		}
	}
	return &processedFileMessage{
		FileName:   filePath,
		StatusCode: Status(StatusSuccessful),
		Message:    "File processed successfully",
		ResultPath: resultPath,
	}
}

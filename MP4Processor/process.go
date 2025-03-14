package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"sync"
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

			filePath := string(msg.Data)
			wg.Add(1)
			go func(filePath string) {
				defer wg.Done()
				publishProcessedFileMessage(filePath, processResultTopic, publish)
			}(filePath)
		}
	}
}

func publishProcessedFileMessage(filePath string, processResultTopic string, publish PublishFunc) {
	PublishStartProcessingMessage(filePath, processResultTopic, publish)

	resultMessage := handleExtractionMessage(filePath)
	byteArray, err := json.Marshal(resultMessage)
	if err != nil {
		log.Fatalf("Error marshaling to JSON:%v\n", err)
		return
	}
	if err := publish(processResultTopic, byteArray); err != nil {
		log.Fatalf("Error publishing message: %v\n", err)
	}
}

func PublishStartProcessingMessage(filePath string, processResultTopic string, publish PublishFunc) {
	message := &processedFileMessage{
		FileName:   filePath,
		StatusCode: Status(StatusProcessing),
		Message:    fmt.Sprintf("Start processsing"),
		ResultPath: "",
	}
	starProcessingMessage, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Error marshaling start processing message to JSON:%v\n", err)
		return
	}
	if err := publish(processResultTopic, starProcessingMessage); err != nil {
		log.Fatalf("Error publishing message: %v\n", err)
	}
}

func handleExtractionMessage(filePath string) *processedFileMessage {
	fmt.Printf("Received file path: %s\n", filePath)

	boxes, err := ExtractInitializationSegment(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return &processedFileMessage{
			FileName:   filePath,
			StatusCode: Status(StatusFailed),
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
			StatusCode: Status(StatusFailed),
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

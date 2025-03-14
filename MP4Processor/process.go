package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"path/filepath"
	"sync"
)

type PublishFunc func(subject string, msg []byte) error

func process(ctx context.Context, msgChan chan *nats.Msg, inputDir, outputDir, processResultTopic string, publish PublishFunc) error {
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

			fileName := string(msg.Data)
			inputPath := filepath.Join(inputDir, fileName)
			outputPath := filepath.Join(outputDir, fileName)

			wg.Add(1)
			go func(filePath string) {
				defer wg.Done()
				err := publishProcessedFileMessage(inputPath, outputPath, processResultTopic, publish)
				if err != nil {
					return
				}
			}(fileName)
		}
	}
}

func publishProcessedFileMessage(filePath, outputPath, processResultTopic string, publish PublishFunc) error {
	resultMessage := handleExtractionMessage(filePath, outputPath)
	byteArray, err := json.Marshal(resultMessage)
	if err != nil {
		log.Fatalf("Error marshaling to JSON:%v\n", err)
		return err
	}
	if err := publish(processResultTopic, byteArray); err != nil {
		log.Fatalf("Error publishing message: %v\n", err)
		return err
	}
	return nil
}

func handleExtractionMessage(filePath, outputPath string) *processedFileMessage {
	fmt.Printf("Received file path: %s\n", filePath)

	boxes, err := ExtractInitializationSegment(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return &processedFileMessage{
			FileName:   filePath,
			StatusCode: StatusFailed,
			Message:    fmt.Sprintf("Failed to process init segment: %v", err),
			ResultPath: "",
		}
	}

	for _, box := range boxes {
		fmt.Printf("MP4Box Type: %s, Size: %d\n", box.Type, box.Size)
	}

	resultPath, err := writeResultIntoFile(outputPath, boxes)

	if err != nil {
		return &processedFileMessage{
			FileName:   filePath,
			StatusCode: StatusFailed,
			Message:    fmt.Sprintf("Failed to write into a file: %v", err),
			ResultPath: "",
		}
	}
	return &processedFileMessage{
		FileName:   filePath,
		StatusCode: StatusSuccessful,
		Message:    "File processed successfully",
		ResultPath: resultPath,
	}
}

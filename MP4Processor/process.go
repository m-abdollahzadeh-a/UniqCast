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
	natsURL             = nats.DefaultURL
	mp4FilePathsTopic   = "mp4FilePaths"
	initialSegmentTopic = "InitialSegmentFilePaths"
	channelBufferSize   = 1024
)

func process(ctx context.Context) error {
	// NATs Connection
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	msgChan := make(chan *nats.Msg, channelBufferSize) // Buffered channel to avoid blocking
	// Subscribe to the NATS subject
	sub, err := nc.ChanSubscribe(mp4FilePathsTopic, msgChan)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := sub.Drain(); err != nil {
			log.Printf("Error draining topic: %v\n", err)
		}
		err := sub.Unsubscribe()
		if err != nil {
			log.Printf("Error unsubscribing topic: %v\n", err)
			return
		}
	}()
	fmt.Println("Subscribed to", mp4FilePathsTopic)

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
				resultMessage := handleMessage(msg)
				byteArray, err := json.Marshal(resultMessage)
				if err != nil {
					fmt.Println("Error marshaling to JSON:", err)
				}
				nc.Publish(initialSegmentTopic, byteArray)
			}(msg)
		}
	}
}

func handleMessage(msg *nats.Msg) *processedFileMessage {
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

func writeResultIntoFile(fileName string, boxes []*MP4Box) (string, error) {
	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	for _, box := range boxes {
		err := writeBox(file, box)
		if err != nil {
			return "", err
		}
	}
	return filepath.Abs(fileName)
}

func writeBox(file *os.File, box *MP4Box) error {
	// Write the Size field (4 bytes)
	if err := binary.Write(file, binary.BigEndian, box.Size); err != nil {
		return err
	}

	typeBytes := []byte(box.Type)
	if len(typeBytes) != 4 {
		panic("Type field must be exactly 4 bytes")
	}
	if _, err := file.Write(typeBytes); err != nil {
		return err
	}

	if _, err := file.Write(box.Data); err != nil {
		return err
	}

	return nil
}

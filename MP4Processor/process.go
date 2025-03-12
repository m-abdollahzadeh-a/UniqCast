package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
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
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	msgChan := make(chan *nats.Msg, channelBufferSize) // Buffered channel to avoid blocking

	// Subscribe to the NATS subject using ChanSubscribe
	sub, err := nc.ChanSubscribe(mp4FilePathsTopic, msgChan)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := sub.Drain(); err != nil {
			log.Printf("Error draining subscription: %v\n", err)
		}
		sub.Unsubscribe()
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
				if err := handleMessage(nc, msg); err != nil {
					log.Printf("Error handling message: %v\n", err)
				}
			}(msg)
		}
	}
}

func handleMessage(nc *nats.Conn, msg *nats.Msg) error {
	filePath := string(msg.Data)
	fmt.Printf("Received file path: %s\n", filePath)

	boxes, err := ExtractInitializationSegment(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	for _, box := range boxes {
		fmt.Printf("MP4Box Type: %s, Size: %d\n", box.Type, box.Size)

	}

	resultPath := writeResultIntoFile("output.mp4", boxes)

	// TODO: do it also for all failed states
	resultMessage := &processedFileMessage{
		FileName:   filePath,
		StatusCode: http.StatusOK,
		Message:    "Successful",
		ResultPath: resultPath,
	}
	byteArray, err := json.Marshal(resultMessage)
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return nil
	}

	nc.Publish(initialSegmentTopic, byteArray)
	return nil
}

func writeResultIntoFile(fileName string, boxes []*MP4Box) string {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	for _, box := range boxes {
		err := writeBox(file, box)
		if err != nil {
			panic(err)
		}
	}
	resultPath, err := filepath.Abs(fileName)
	return resultPath
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

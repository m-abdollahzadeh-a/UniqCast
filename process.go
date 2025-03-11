package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

func process() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	subject := "mp4FilePaths"

	msgChan := make(chan *nats.Msg, 1024) // Buffered channel to avoid blocking

	// Subscribe to the NATS subject using ChanSubscribe
	sub, err := nc.ChanSubscribe(subject, msgChan)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := sub.Drain(); err != nil {
			log.Printf("Error draining subscription: %v\n", err)
		}
		sub.Unsubscribe()
	}()

	fmt.Println("Subscribed to", subject)

	for msg := range msgChan {
		filePath := string(msg.Data)
		fmt.Printf("Received file path: %s\n", filePath)

		boxes, err := ExtractInitializationSegment(filePath)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		for _, box := range boxes {
			fmt.Printf("MP4Box Type: %s, Size: %d\n", box.Type, box.Size)

		}

		err = writeResultIntoFile("output.mp4", boxes)
		if err != nil {
			return
		}
	}
}

func writeResultIntoFile(fileName string, boxes []*MP4Box) error {
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
	return nil
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

package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

func process() {
	// Connect to the NATS server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Define the subject to subscribe to
	subject := "mp4FilePaths"

	// Create a channel to receive NATS messages
	msgChan := make(chan *nats.Msg, 1024) // Buffered channel to avoid blocking

	// Subscribe to the NATS subject using ChanSubscribe
	sub, err := nc.ChanSubscribe(subject, msgChan)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		// Drain the subscription before unsubscribing
		if err := sub.Drain(); err != nil {
			log.Printf("Error draining subscription: %v\n", err)
		}
		// Unsubscribe after draining
		sub.Unsubscribe()
	}()

	// Keep the program running to receive messages
	fmt.Println("Subscribed to", subject)

	for msg := range msgChan {
		filePath := string(msg.Data)
		fmt.Printf("Received file path: %s\n", filePath)

		// Process the file
		boxes, err := ExtractInitializationSegment(filePath)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		for _, box := range boxes {
			fmt.Printf("MP4Box Type: %s, Size: %d\n", box.Type, box.Size)
		}
	}
}

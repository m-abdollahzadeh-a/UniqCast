package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestPublishProcessedFileMessage(t *testing.T) {
	// Start NATS container using testcontainers
	ctx := context.Background()
	natsContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "nats:2.9.22",
			ExposedPorts: []string{"4222/tcp"},
			WaitingFor:   wait.ForLog("Listening for client connections on 0.0.0.0:4222"),
		},
		Started: true,
	})
	if err != nil {
		t.Fatalf("Failed to start NATS container: %v", err)
	}
	defer natsContainer.Terminate(ctx)

	// Get NATS connection URL
	natsHost, err := natsContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get NATS container host: %v", err)
	}
	natsPort, err := natsContainer.MappedPort(ctx, "4222")
	if err != nil {
		t.Fatalf("Failed to get NATS container port: %v", err)
	}
	natsURL := fmt.Sprintf("nats://%s:%s", natsHost, natsPort.Port())

	// Connect to NATS server
	nc, err := nats.Connect(natsURL)
	if err != nil {
		t.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	publishFunc := func(topic string, message []byte) error {
		return nc.Publish(topic, message)
	}

	// Subscribe to the processResultTopic to verify messages
	processResultTopic := "test.topic"
	receivedMessages := make(chan *processedFileMessage, 1)
	subscription, err := nc.Subscribe(processResultTopic, func(msg *nats.Msg) {
		var result processedFileMessage
		if err := json.Unmarshal(msg.Data, &result); err != nil {
			t.Errorf("Failed to unmarshal message: %v", err)
			return
		}
		receivedMessages <- &result
	})
	if err != nil {
		t.Fatalf("Failed to subscribe to topic: %v", err)
	}
	defer subscription.Unsubscribe()

	// Test cases
	tests := []struct {
		name          string
		filePath      string
		outputPath    string
		expectedError bool
		expectedMsg   *processedFileMessage
	}{
		{
			name:          "successful processing",
			filePath:      "../video.mp4",
			outputPath:    "./video.mp4",
			expectedError: false,
			expectedMsg: &processedFileMessage{
				FileName:   "../video.mp4",
				StatusCode: StatusSuccessful,
				Message:    "File processed successfully",
				ResultPath: "./video.mp4",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := publishProcessedFileMessage(tt.filePath, tt.outputPath, processResultTopic, publishFunc)
			if (err != nil) != tt.expectedError {
				t.Errorf("publishProcessedFileMessage() error = %v, expectedError %v", err, tt.expectedError)
			}

			// Wait for the message to be received
			select {
			case msg := <-receivedMessages:
				if msg.FileName != tt.expectedMsg.FileName ||
					msg.StatusCode != tt.expectedMsg.StatusCode ||
					msg.Message != tt.expectedMsg.Message {
					t.Errorf("Received message = %v, expected %v", msg, tt.expectedMsg)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("Timeout waiting for message")
			}
		})
	}
}

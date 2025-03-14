package main

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"testing"
	"time"

	"MP4Processor/model"
	"MP4Processor/processor"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestProcessNatsMessage(t *testing.T) {
	// Test cases
	tests := []struct {
		name       string
		filePath   string
		outputPath string
		res        model.ProcessedFileMessage
	}{
		{
			name:       "successful processing",
			filePath:   "video.mp4",
			outputPath: "/tmp",
			res: model.ProcessedFileMessage{
				FileName:   "video.mp4",
				StatusCode: model.StatusSuccessful,
				Message:    "File processed successfully",
				ResultPath: "/tmp",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup processor
			p := processor.New(tt.outputPath)

			natsRes := processNatsMessage(p, tt.filePath, tt.outputPath)
			require.Equal(t, tt.res, natsRes, "res")
		})
	}
}

func TestProcessNatsMessage_E2E(t *testing.T) {
	// Start a NATS server using Testcontainers
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "nats:2.9.22",
		ExposedPorts: []string{"4222/tcp"},
		WaitingFor:   wait.ForLog("Listening for client connections on 0.0.0.0:4222"),
	}
	natsContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	defer natsContainer.Terminate(ctx)

	// Get the NATS server URL
	natsURL, err := natsContainer.Endpoint(ctx, "")
	require.NoError(t, err)

	// Connect to the NATS server
	nc, err := nats.Connect(natsURL)
	require.NoError(t, err)
	defer nc.Close()

	// Test cases
	tests := []struct {
		name       string
		filePath   string
		outputPath string
		res        model.ProcessedFileMessage
	}{
		{
			name:       "successful processing",
			filePath:   "video.mp4",
			outputPath: "/tmp",
			res: model.ProcessedFileMessage{
				FileName:   "video.mp4",
				StatusCode: model.StatusSuccessful,
				Message:    "File processed successfully",
				ResultPath: "/tmp",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup processor
			p := processor.New(tt.outputPath)

			// Process the message
			natsRes := processNatsMessage(p, tt.filePath, tt.outputPath)

			// Serialize the result to JSON
			resBytes, err := json.Marshal(natsRes)
			require.NoError(t, err)

			// Subscribe to the NATS subject to verify the result
			sub, err := nc.SubscribeSync("processed.file")
			require.NoError(t, err)

			// Add a small delay to ensure the subscription is active
			time.Sleep(500 * time.Millisecond)

			// Publish the result to NATS
			err = nc.Publish("processed.file", resBytes)
			require.NoError(t, err)

			// Wait for the message
			msg, err := sub.NextMsg(10 * time.Second)
			require.NoError(t, err)

			// Unmarshal the message
			var receivedRes model.ProcessedFileMessage
			err = json.Unmarshal(msg.Data, &receivedRes)
			require.NoError(t, err)

			// Verify the result
			require.Equal(t, tt.res, receivedRes, "res")
		})
	}
}

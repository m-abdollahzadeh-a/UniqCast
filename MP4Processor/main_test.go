package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"MP4Processor/model"
	"MP4Processor/processor"
)

func TestProcessNatsMessage(t *testing.T) {
	// Test cases
	tests := []struct {
		name       string
		filePath   string
		outputPath string
		res        model.ProcessedFileMessage
	}{
		// todo check filepath for input and output
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

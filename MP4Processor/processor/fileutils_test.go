package processor

import (
	"encoding/binary"
	"os"
	"testing"

	"MP4Processor/model"
)

func TestWriteBox(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFileName := tmpFile.Name()
	defer os.Remove(tmpFileName)

	// Define test data
	box := &model.MP4Box{
		Size: 8,
		Type: "ftyp",
		Data: []byte{0x00, 0x01, 0x02, 0x03},
	}

	if err := writeBox(tmpFile, box); err != nil {
		t.Fatalf("writeBox failed: %v", err)
	}

	// Read the file back to verify its contents
	file, err := os.Open(tmpFileName)
	if err != nil {
		t.Fatalf("Failed to open file for verification: %v", err)
	}
	defer file.Close()

	// Verify the box
	var size uint32
	if err := binary.Read(file, binary.BigEndian, &size); err != nil {
		t.Fatalf("Failed to read size: %v", err)
	}
	if size != box.Size {
		t.Errorf("Expected size %d, got %d", box.Size, size)
	}

	typeBytes := make([]byte, 4)
	if _, err := file.Read(typeBytes); err != nil {
		t.Fatalf("Failed to read type: %v", err)
	}
	if string(typeBytes) != box.Type {
		t.Errorf("Expected type %s, got %s", box.Type, string(typeBytes))
	}

	data := make([]byte, len(box.Data))
	if _, err := file.Read(data); err != nil {
		t.Fatalf("Failed to read data: %v", err)
	}
	if string(data) != string(box.Data) {
		t.Errorf("Expected data %v, got %v", box.Data, data)
	}
}

package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
)

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
		return fmt.Errorf("type field must be exactly 4 bytes")
	}
	if _, err := file.Write(typeBytes); err != nil {
		return err
	}

	if _, err := file.Write(box.Data); err != nil {
		return err
	}

	return nil
}

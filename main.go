package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

func ReadBox(file *os.File) (*MP4Box, error) {
	var size uint32
	err := binary.Read(file, binary.BigEndian, &size)
	if err != nil {
		return nil, err
	}

	typeBytes := make([]byte, 4)
	_, err = file.Read(typeBytes)
	if err != nil {
		return nil, err
	}
	boxType := string(typeBytes)

	dataSize := size - 8 // Subtract 8 bytes for size and type fields
	data := make([]byte, dataSize)
	_, err = file.Read(data)
	if err != nil {
		return nil, err
	}

	return &MP4Box{
		Size: size,
		Type: boxType,
		Data: data,
	}, nil
}

func ExtractInitializationSegment(filePath string) ([]*MP4Box, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	var boxes []*MP4Box

	for {
		box, err := ReadBox(file)
		if err != nil {
			return nil, err
		}

		boxes = append(boxes, box)

		if box.Type == "moov" {
			break
		}
	}

	return boxes, nil
}

func main() {
	filePath := "../video.mp4"
	boxes, err := ExtractInitializationSegment(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, box := range boxes {
		fmt.Printf("MP4Box Type: %s, Size: %d\n", box.Type, box.Size)
	}
}

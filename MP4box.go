package main

import (
	"encoding/binary"
	"log"
	"os"
)

type MP4Box struct {
	Size uint32 // Size of the box including header and data
	Type string
	Data []byte
}

func ReadBox(file *os.File) (*MP4Box, error) {
	var size uint32

	err := binary.Read(file, binary.BigEndian, &size)
	if err != nil {
		log.Printf("Failed to read box size: %v", err)
		return nil, err
	}

	typeBytes := make([]byte, 4)

	_, err = file.Read(typeBytes)
	if err != nil {
		log.Printf("Failed to read box type: %v", err)
		return nil, err
	}
	boxType := string(typeBytes)

	dataSize := size - 8 // Subtract 8 bytes for size and type fields
	data := make([]byte, dataSize)
	_, err = file.Read(data)
	if err != nil {
		log.Printf("Failed to read box data: %v", err)
		return nil, err
	}

	log.Printf("Successfully read box: Type=%s, Size=%d", boxType, size)
	return &MP4Box{
		Size: size,
		Type: boxType,
		Data: data,
	}, nil
}

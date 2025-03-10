package main

import (
	"encoding/binary"
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

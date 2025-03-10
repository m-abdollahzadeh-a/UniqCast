package main

import "os"

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

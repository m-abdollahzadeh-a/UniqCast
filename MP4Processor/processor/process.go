package processor

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"MP4Processor/model"
)

type Processor struct {
	outputPath string
}

func New(outputPath string) *Processor {
	return &Processor{outputPath: outputPath}
}

func (p *Processor) ProcessMP4(inputFile string) (err error) {
	filename := filepath.Base(inputFile)
	outputFile := filepath.Join(p.outputPath, filename)

	// Extract init segment (initial boxes)
	initBoxes, err := extractInitSegment(inputFile)
	if err != nil {
		return err
	}

	// Write init segment to outputFile
	return writeResultIntoFile(outputFile, initBoxes)
}

func extractInitSegment(filePath string) (boxes []*model.MP4Box, err error) {
	fmt.Printf("Received file path: %s\n", filePath)

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	// Close file
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("falied to close file %s: %s", filePath, err)
		}
	}(file)

	// Find init boxes
	for {
		box, err := readBox(file)
		if err != nil {
			log.Printf("Failed to read box: %v", err)
			return nil, err
		}

		boxes = append(boxes, box)
		if box.Type == model.EndOfInitialSegmentType {
			log.Printf("Found moov box, stopping extraction")
			break
		}
	}

	for _, box := range boxes {
		fmt.Printf("MP4Box Type: %s, Size: %d\n", box.Type, box.Size)
	}

	return boxes, nil
}

func readBox(file *os.File) (box *model.MP4Box, err error) {
	var size uint32

	err = binary.Read(file, binary.BigEndian, &size)
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
	box = &model.MP4Box{
		Size: size,
		Type: boxType,
		Data: data,
	}
	return box, nil
}

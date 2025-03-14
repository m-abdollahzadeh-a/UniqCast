package processor

import (
	"encoding/binary"
	"fmt"
	"os"

	"MP4Processor/model"
)

func writeResultIntoFile(fileName string, boxes []*model.MP4Box) (err error) {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, box := range boxes {
		err := writeBox(file, box)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeBox(file *os.File, box *model.MP4Box) error {
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

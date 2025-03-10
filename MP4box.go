package main

type MP4Box struct {
	Size uint32 // Size of the box including header and data
	Type string
	Data []byte
}

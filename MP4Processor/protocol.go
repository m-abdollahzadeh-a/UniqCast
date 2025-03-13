package main

type Status string

const (
	StatusSuccessful         string = "successful"
	StatusFailedProcessing   string = "failed_processing"
	StatusFailedToWritInFile string = "failed_to_write_in_file"
)

type processedFileMessage struct {
	FileName   string `json:"file_name"`
	StatusCode Status `json:"status_code"`
	Message    string `json:"Message"`
	ResultPath string `json:"result_path"`
}

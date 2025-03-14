package main

type Status string

const (
	StatusSuccessful Status = "Successful"
	StatusFailed     Status = "Failed"
)

type processedFileMessage struct {
	FileName   string `json:"file_name"`
	StatusCode Status `json:"status_code"`
	Message    string `json:"Message"`
	ResultPath string `json:"result_path"`
}

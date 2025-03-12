package main

type processedFileMessage struct {
	FileName   string `json:"file_name"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"Message"`
	ResultPath string `json:"result_path"`
}

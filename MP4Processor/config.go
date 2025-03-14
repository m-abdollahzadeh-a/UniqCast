package main

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	NatsURL            string `yaml:"natsURL"`
	Mp4FilePathsTopic  string `yaml:"mp4FilePathsTopic"`
	ProcessResultTopic string `yaml:"processResultTopic"`
}

func loadConfig(filePath string) Config {
	var config Config

	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML: %v", err)
	}
	return config
}

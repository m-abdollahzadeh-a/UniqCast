package main

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	NATS struct {
		URL                string `yaml:"url"`
		Mp4FilePathsTopic  string `yaml:"mp4FilePathsTopic"`
		ProcessResultTopic string `yaml:"processResultTopic"`
	} `yaml:"nats"`

	File struct {
		InputPath  string `yaml:"inputPath"`
		OutputPath string `yaml:"outputPath"`
	} `yaml:"file"`
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

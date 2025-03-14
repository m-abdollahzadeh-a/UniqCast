package config

import (
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	NATS NATS
	File File
}

type NATS struct {
	URL                string
	Mp4FilePathsTopic  string
	ProcessResultTopic string
	BufferSize         int
}

type File struct {
	InputPath  string
	OutputPath string
}

func LoadConfig() (*Config, error) {
	fs := pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError) // flag set
	fs.String("nats-url", "nats://127.0.0.1:4222", "nats url")
	fs.String("nats-file-path-topic", "mp4FilePaths", "nats file path topic")
	fs.String("nats-process-result-topic", "InitialSegmentFilePaths", "nats process result topic")
	fs.Int("nats-buffer-size", 1024, "nats channel buffer size")
	fs.String("file-input-path", "/tmp/inputs/", "file input path")
	fs.String("file-output-path", "/tmp/outputs/", "file output path")
	if err := fs.Parse(os.Args[1:]); err != nil {
		// printUsage(f)
		return nil, err
	}

	v := viper.New()
	if err := v.BindPFlags(fs); err != nil {
		// printUsage(f)
		return nil, err
	}

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	conf := &Config{
		NATS: loadNats(v),
		File: loadFile(v),
	}

	return conf, nil
}

func loadNats(v *viper.Viper) NATS {
	return NATS{
		URL:                v.GetString("nats-url"),
		Mp4FilePathsTopic:  v.GetString("nats-file-path-topic"),
		ProcessResultTopic: v.GetString("nats-process-result-topic"),
		BufferSize:         v.GetInt("nats-buffer-size"),
	}
}

func loadFile(v *viper.Viper) File {
	return File{
		InputPath:  v.GetString("file-input-path"),
		OutputPath: v.GetString("file-output-path"),
	}
}

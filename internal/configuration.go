package internal

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	LogLevel         string `yaml:"log_level"`
	AccountBatchSize int    `yaml:"account_batch_size"`
}

type Account struct {
	Account string `json:"account"`
	Region  string `json:"region"`
}

func ReadConfig(configPath string) (*Config, error) {
	filePath := filepath.Join(configPath, "config.yaml")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read config file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("unable to parse config file: %w", err)
	}

	return &config, nil
}

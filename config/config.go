package config

import (
	"contract-reader/utils"
	"encoding/json"
	"github.com/KyberNetwork/logger"
	"github.com/Netflix/go-env"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	ContractPath string `env:"CONTRACT_PATH" default:"./config/contract.json"` //file mount to docker
	Topic        string `env:"TOPIC" default:"test"`
	ChainId      int64  `env:"CHAIN_ID" default:"1"`
}

func NewConfig() *Config {
	var config Config
	_, err := env.UnmarshalFromEnviron(&config)
	if err != nil {
		logger.Fatalf("failed to unmarshal config from env: %v", err)
	}
	return &config
}

func (c *Config) GetListContracts() []utils.Contract {
	_, b, _, _ := runtime.Caller(0) //nolint:dogsled
	basePath := filepath.Dir(b)
	path := filepath.Join(basePath, c.ContractPath)
	file, err := os.Open(path)
	if err != nil {
		logger.Fatalf("failed to open config file: %v", err)
	}
	defer file.Close()
	jsonStruct := struct {
		Contracts []utils.Contract `json:"contracts"`
	}{}
	value, err := io.ReadAll(file)
	if err != nil {
		logger.Fatalf("failed to read config file: %v", err)
	}
	err = json.Unmarshal(value, &jsonStruct)
	if err != nil {
		logger.Fatalf("failed to unmarshal config file: %v", err)
	}
	return jsonStruct.Contracts
}

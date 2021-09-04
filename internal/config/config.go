package config

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
)

type ServerConfig struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
}

type DBConfig struct {
	Host            string `json:"host"`
	Port            uint16 `json:"port"`
	DbName          string `json:"db_name"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	MaxConnections  uint32 `json:"max_connections"`

	// The application at the startup tries to connect to a DB ConnectionTries times
	// waiting RetryPeriodMs milliseconds between retries
	ConnectionTries uint32 `json:"connection_tries"`
	RetryPeriodMs   uint32 `json:"reconnect_period_ms"`
}

type SettingsConfig struct {
	RepoFlushBatchSize uint32 `json:"repo_flush_batch_size"`
	InternalBufferSize uint32 `json:"internal_buffer_size"`
	FlushPeriodMs      uint32 `json:"flush_period_ms"`
}

type ApplicationConfig struct {
	Server   ServerConfig   `json:"server_config"`
	Db       DBConfig       `json:"db_config"`
	Settings SettingsConfig `json:"settings_config"`
}

func ReadApplicationConfig(path string) (*ApplicationConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer closeFile(file)

	decoder := json.NewDecoder(file)
	config := &ApplicationConfig{}
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}

func closeFile(file *os.File) {
	if err := file.Close(); err != nil {
		log.Error().
			Str("reason", "unable to close application config file due to an error").
			Msgf("%v", err)
	}
}

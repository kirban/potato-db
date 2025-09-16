package config

import (
	"bytes"
	"errors"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"slices"
)

var (
	ValidLogLevels   = []string{"debug", "info", "warn", "error", "panic", "fatal"}
	ValidLogOutputs  = []string{"stdout", "stderr"}
	ValidEngineTypes = []string{"in-memory", "disk"}
)

type Configurable[T any] interface {
	ParseYaml(path string) (*T, error)
	ValidateConfig() error
}

type Config struct {
	app       *AppConfigOptions    `yaml:"app"`
	tcpServer *ServerConfigOptions `yaml:"tcp_server"`
	db        *DbConfigOptions     `yaml:"db"`
}

type AppConfigOptions struct {
	LogLevel  string `yaml:"level"`
	LogOutput string `yaml:"output"`
}

type DbConfigOptions struct {
	EngineType string `yaml:"engine_type"`
}

type ServerConfigOptions struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

var ServerConfigDefaults = &ServerConfigOptions{
	Host: "127.0.0.1",
	Port: 8282,
}

var DbConfigDefaults = &DbConfigOptions{
	EngineType: "in-memory",
}

var AppConfigDefaults = &AppConfigOptions{
	LogLevel:  "info",
	LogOutput: "stdout",
}

func (c *Config) ParseYaml(configPath string) (*Config, error) {
	file, err := os.ReadFile(configPath)

	if err != nil {
		return nil, errors.New("failed to read config file: " + err.Error())
	}

	reader := bytes.NewReader(file)
	data, err := io.ReadAll(reader)

	if err != nil {
		return nil, errors.New("failed to read buffer: " + err.Error())
	}

	var config Config
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, errors.New("failed to parse config file: " + err.Error())
	}

	if err = config.ValidateConfig(); err != nil {
		return nil, errors.New("failed to validate config file: " + err.Error())
	}

	return &config, nil
}

func (c *Config) ValidateConfig() error {
	if c.app == nil {
		return errors.New("app is required")
	} else {
		if c.app.LogLevel == "" {
			c.app.LogLevel = AppConfigDefaults.LogLevel
		} else if !slices.Contains(ValidLogLevels, c.app.LogLevel) {
			return errors.New("invalid log level")
		}

		if c.app.LogOutput == "" {
			c.app.LogOutput = AppConfigDefaults.LogOutput
		} else if !slices.Contains(ValidLogOutputs, c.app.LogOutput) {
			return errors.New("invalid log output")
		}
	}

	if c.tcpServer == nil {
		return errors.New("tcp_server is required")
	} else {
		if c.tcpServer.Host == "" {
			c.tcpServer.Host = ServerConfigDefaults.Host
		}

		if c.tcpServer.Port == 0 {
			c.tcpServer.Port = ServerConfigDefaults.Port
		} else if c.tcpServer.Port < 1024 || c.tcpServer.Port > 65535 {
			return errors.New("invalid tcp server port")
		}
	}

	if c.db == nil {
		return errors.New("db is required")
	} else {
		if c.db.EngineType == "" {
			c.db.EngineType = DbConfigDefaults.EngineType
		} else if !slices.Contains(ValidEngineTypes, c.db.EngineType) {
			return errors.New("invalid db engine type")
		}
	}

	return nil
}

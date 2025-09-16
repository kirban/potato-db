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
	parseYaml(path string) (*T, error)
	validateConfig() error
}

type Config struct {
	App       *AppConfigOptions    `yaml:"app"`
	TcpServer *ServerConfigOptions `yaml:"tcp_server"`
	Db        *DbConfigOptions     `yaml:"db"`
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

func (c *Config) parseYaml(rawData []byte, config *Config) (*Config, error) {
	if err := yaml.Unmarshal(rawData, &config); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) validateConfig() error {
	if c.App == nil {
		return errors.New("App is required")
	} else {
		if c.App.LogLevel == "" {
			c.App.LogLevel = AppConfigDefaults.LogLevel
		} else if !slices.Contains(ValidLogLevels, c.App.LogLevel) {
			return errors.New("invalid log level")
		}

		if c.App.LogOutput == "" {
			c.App.LogOutput = AppConfigDefaults.LogOutput
		} else if !slices.Contains(ValidLogOutputs, c.App.LogOutput) {
			return errors.New("invalid log output")
		}
	}

	if c.TcpServer == nil {
		return errors.New("tcp_server is required")
	} else {
		if c.TcpServer.Host == "" {
			c.TcpServer.Host = ServerConfigDefaults.Host
		}

		if c.TcpServer.Port == 0 {
			c.TcpServer.Port = ServerConfigDefaults.Port
		} else if c.TcpServer.Port < 1024 || c.TcpServer.Port > 65535 {
			return errors.New("invalid tcp server port")
		}
	}

	if c.Db == nil {
		return errors.New("Db is required")
	} else {
		if c.Db.EngineType == "" {
			c.Db.EngineType = DbConfigDefaults.EngineType
		} else if !slices.Contains(ValidEngineTypes, c.Db.EngineType) {
			return errors.New("invalid Db engine type")
		}
	}

	return nil
}

func NewConfig(configPath string) (*Config, error) {
	// load from file
	file, err := os.ReadFile(configPath)

	if err != nil {
		return nil, errors.New("failed to read config file: " + err.Error())
	}

	reader := bytes.NewReader(file)
	data, err := io.ReadAll(reader)

	if err != nil {
		return nil, errors.New("failed to read buffer: " + err.Error())
	}

	cfg := &Config{
		App:       &AppConfigOptions{},
		TcpServer: &ServerConfigOptions{},
		Db:        &DbConfigOptions{},
	}

	// parse
	cfg, err = cfg.parseYaml(data, cfg)

	if err != nil {
		return nil, errors.New("failed to parse config file: " + err.Error())
	}

	// validate
	if err = cfg.validateConfig(); err != nil {
		return nil, errors.New("failed to validate config file: " + err.Error())
	}

	return cfg, nil
}

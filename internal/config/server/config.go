package server

import (
	"embed"
	"fmt"
	"log/slog"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const defaultConfigFile = "default.yaml"

//go:embed *.yaml
var defaultConfig embed.FS

type EngineType string

const (
	EngineTypeInMemory EngineType = "in_memory"
)

type Config struct {
	Engine struct {
		Type EngineType `yaml:"type"`
	} `yaml:"engine"`
	Network struct {
		Address        string        `yaml:"address"`
		MaxConnections int           `yaml:"max_connections"`
		MaxBufferSize  int           `yaml:"max_buffer_size_bytes"`
		IdleTimeout    time.Duration `yaml:"idle_timeout_ms"`
	} `yaml:"network"`
	Logging struct {
		Level slog.Level `yaml:"level"`
	} `yaml:"logging"`
}

func readDefaultConfig() (*Config, error) {
	config := &Config{}

	data, err := defaultConfig.ReadFile(defaultConfigFile)
	if err != nil {
		return nil, fmt.Errorf("failed read default config: %w", err)
	}

	if err = yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed parse default config: %w", err)
	}

	return config, nil
}

func readConfigWithOverride(config *Config, fileName string) (*Config, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed read config %q: %w", fileName, err)
	}

	if err = yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed parse config %q: %w", fileName, err)
	}

	return config, nil
}

func Load() (*Config, error) {
	config, err := readDefaultConfig()
	if err != nil {
		return nil, err
	}

	flags := readFlags(flagValues{
		address: config.Network.Address,
		verbose: false,
	})

	if flags.config != "" {
		config, err = readConfigWithOverride(config, flags.config)
		if err != nil {
			return nil, err
		}
	}

	config.Network.Address = flags.address

	if flags.verbose {
		config.Logging.Level = slog.LevelDebug
	}

	return config, nil
}

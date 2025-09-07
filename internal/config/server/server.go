package server

import (
	"flag"
	"log/slog"
)

type Config struct {
	Addr string `yaml:"addr"`

	LogLevel slog.Level `yaml:"log_level"`
}

func LoadConfig() (Config, error) {
	server := getDefaultConfig()

	var verbose bool
	flag.BoolVar(&verbose, "verbose", verbose, "Verbosity")
	flag.StringVar(&server.Addr, "address", server.Addr, "Server address")

	flag.Parse()

	if verbose {
		server.LogLevel = slog.LevelDebug
	}

	return server, nil
}

func getDefaultConfig() Config {
	return Config{
		Addr:     ":8081",
		LogLevel: slog.LevelError,
	}
}

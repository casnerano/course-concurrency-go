package config

import (
	"flag"
	"log/slog"
)

var (
	ServerName    = "memdb"
	ServerVersion = "dev"
)

type Server struct {
	Addr string `yaml:"addr"`

	LogLevel slog.Level `yaml:"log_level"`
}

var server = Server{
	Addr:     ":8081",
	LogLevel: slog.LevelError,
}

func InitServer() error {
	var verbose bool
	flag.BoolVar(&verbose, "verbose", verbose, "Verbosity")
	flag.StringVar(&server.Addr, "address", server.Addr, "Server address")

	flag.Parse()

	if verbose {
		server.LogLevel = slog.LevelDebug
	}

	return nil
}

func GetServer() Server {
	return server
}

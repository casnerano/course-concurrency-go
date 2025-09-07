package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/casnerano/course-concurrency-go/internal/config"
	"github.com/casnerano/course-concurrency-go/internal/logger"
)

func main() {
	if err := config.InitServer(); err != nil {
		panic(err.Error())
	}

	cfg := config.GetServer()

	logger.Init(config.ServerName, cfg.LogLevel)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	logger.Info("Stopped server..")
}

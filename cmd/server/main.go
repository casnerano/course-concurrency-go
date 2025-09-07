package main

import (
	"context"
	"os/signal"
	"syscall"

	serverconfig "github.com/casnerano/course-concurrency-go/internal/config/server"
	"github.com/casnerano/course-concurrency-go/internal/logger"
)

func main() {
	config, configErr := serverconfig.LoadConfig()
	if configErr != nil {
		panic(configErr.Error())
	}

	logger.Init("memdb", config.LogLevel)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	_ = getDatabase()

	<-ctx.Done()

	logger.Info("Stopped server..")
}

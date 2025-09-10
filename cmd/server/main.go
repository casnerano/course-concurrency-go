package main

import (
	"context"
	"os/signal"
	"syscall"

	config "github.com/casnerano/course-concurrency-go/internal/config/server"
	"github.com/casnerano/course-concurrency-go/internal/database"
	"github.com/casnerano/course-concurrency-go/internal/database/compute"
	"github.com/casnerano/course-concurrency-go/internal/database/storage"
	"github.com/casnerano/course-concurrency-go/internal/database/storage/engine/memory"
	"github.com/casnerano/course-concurrency-go/internal/logger"
	"github.com/casnerano/course-concurrency-go/internal/network"
	"github.com/casnerano/course-concurrency-go/internal/network/protocol"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err.Error())
	}

	logger.Init("memdb", cfg.LogLevel)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err = runServer(ctx, cfg.Addr); err != nil {
		logger.Error("server internal error: " + err.Error())
	}

	<-ctx.Done()

	logger.Info("stopped server..")
}

func runServer(ctx context.Context, addr string) error {
	server := network.NewServer(
		addr,
		protocol.NewJSON(),
		database.New(
			compute.New(),
			storage.New(
				memory.New(),
			),
		),
	)

	return server.Start(ctx)
}

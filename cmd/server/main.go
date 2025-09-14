package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	config "github.com/casnerano/course-concurrency-go/internal/config/server"
	"github.com/casnerano/course-concurrency-go/internal/database"
	"github.com/casnerano/course-concurrency-go/internal/database/compute"
	"github.com/casnerano/course-concurrency-go/internal/database/storage"
	memory_engine "github.com/casnerano/course-concurrency-go/internal/database/storage/engine/memory"
	"github.com/casnerano/course-concurrency-go/internal/logger"
	"github.com/casnerano/course-concurrency-go/internal/network"
	"github.com/casnerano/course-concurrency-go/internal/network/protocol"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err.Error())
	}

	logger.Init("memdb", cfg.Logging.Level)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var server *network.Server
	server, err = buildServer(cfg)
	if err != nil {
		logger.Error("server configuration error: " + err.Error())
	}

	if err = server.Start(ctx); err != nil {
		logger.Error("server internal error: " + err.Error())
	}

	<-ctx.Done()

	logger.Info("stopped server..")
}

func buildServer(cfg *config.Config) (*network.Server, error) {
	var engine storage.Engine
	switch engineType := cfg.Engine.Type; engineType {
	case config.EngineTypeInMemory:
		engine = memory_engine.New()
	default:
		return nil, fmt.Errorf("unknown engine type: %s", engineType)
	}

	networkOptions := network.ServerOptions{
		Address:        cfg.Network.Address,
		MaxConnections: cfg.Network.MaxConnections,
		IdleTimeout:    cfg.Network.IdleTimeout,
	}

	server := network.NewServer(
		protocol.NewJSON(
			protocol.WithMaxBufferSize(
				cfg.Network.MaxBufferSize,
			),
		),
		database.New(
			compute.New(),
			storage.New(engine),
		),
		networkOptions,
	)

	return server, nil
}

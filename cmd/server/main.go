package main

import (
	"context"
	"github.com/kirban/potato-db/internal/config"
	"github.com/kirban/potato-db/internal/network"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		configPath = "config.potato.yaml"
	}

	cfg, err := config.NewConfig(configPath)

	if err != nil {
		log.Fatal(err)
	}

	logger, err := zap.NewDevelopment() // todo вынести инициализацию логера

	server, err := network.NewTCPServer(logger, cfg.TcpServer)

	if err != nil {
		logger.Fatal("failed to create server", zap.Error(err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		err := server.StartAndServe()
		if err != nil {
			logger.Fatal("failed starting server", zap.Error(err))
		}
	}()

	logger.Info("Server started. Press CTRL+C to stop")
	<-ctx.Done()
	logger.Info("Got exit signal. Gracefully shutdown.")

}

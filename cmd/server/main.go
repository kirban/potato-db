package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kirban/potato-db/internal/config"
	loggerModule "github.com/kirban/potato-db/internal/logger"
	"github.com/kirban/potato-db/internal/network"
	"go.uber.org/zap"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		configPath = "config.potato.yaml"
	}

	cfg, err := config.NewConfig(configPath)

	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger, err := loggerModule.NewLogger(cfg)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := logger.Sync(); err != nil {
			// ожидаемые ошибки stdout/stderr на macOS можно игнорировать
			log.Printf("logger sync error: %v", err)
		}
	}()

	defer func() {
		if r := recover(); r != nil {
			logger.Panic("uncaught panic", zap.Any("recovery", r))
		}
	}()

	server, err := network.NewTCPServer(logger, cfg.TcpServer)

	if err != nil {
		logger.Fatal("failed to create server", zap.Error(err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.StartAndServe(ctx); err != nil {
			logger.Fatal("failed starting server", zap.Error(err))
		}
	}()

	logger.Info("Server started. Press CTRL+C to stop")
	<-ctx.Done()
	logger.Info("Got exit signal. Gracefully shutdown.")

}

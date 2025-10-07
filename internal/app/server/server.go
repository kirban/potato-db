package server

import (
	"context"
	"github.com/kirban/potato-db/internal/config"
	"github.com/kirban/potato-db/internal/db"
	loggerModule "github.com/kirban/potato-db/internal/logger"
	"github.com/kirban/potato-db/internal/network"
	"github.com/kirban/potato-db/internal/network/handlers"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const DefaultConfigPath = "config.potato.yaml"

type AppServer struct {
	config *config.Config
	logger *zap.Logger
	db     *db.Database
	server *network.TCPServer
}

func NewAppServer() (*AppServer, error) {
	app := &AppServer{}

	if err := app.initDeps(); err != nil {
		return nil, err
	}

	return app, nil
}

func (s *AppServer) Run() {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Panic("uncaught panic", zap.Any("recovery", r))
		}
	}()

	defer func() {
		if err := s.logger.Sync(); err != nil {
			// ожидаемые ошибки stdout/stderr на macOS можно игнорировать
			log.Printf("logger sync error: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := s.server.StartAndServe(ctx); err != nil {
			s.logger.Fatal("failed starting server", zap.Error(err))
		}
	}()

	s.logger.Info("Server started. Press CTRL+C to stop")
	<-ctx.Done()
	s.logger.Info("Got exit signal. Gracefully shutdown.")
}

func (s *AppServer) initDeps() error {
	deps := []func() error{
		s.initConfig,
		s.initLogger,
		s.initDatabase,
		s.initServer,
	}

	for _, dep := range deps {
		if err := dep(); err != nil {
			return err
		}
	}

	return nil
}

func (s *AppServer) initConfig() error {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		configPath = DefaultConfigPath
	}

	cfg, err := config.NewConfig(configPath)

	if err != nil {
		return err
	}

	s.config = cfg
	return nil
}

func (s *AppServer) initLogger() error {
	logger, err := loggerModule.NewLogger(s.config)

	if err != nil {
		log.Fatal(err)
	}

	s.logger = logger
	return nil
}

func (s *AppServer) initDatabase() error {
	database := db.NewDbBuilder(s.logger).
		InitStorage().
		InitCompute().
		Build()

	s.db = database
	return nil
}

func (s *AppServer) initServer() error {
	handler := &handlers.DatabaseHandler{
		Db: s.db,
	}

	server, err := network.NewTCPServer(s.logger, s.config.TcpServer, handler)

	if err != nil {
		s.logger.Fatal("failed to create server", zap.Error(err))
	}

	s.server = server
	return nil
}

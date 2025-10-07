package cli

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/kirban/potato-db/internal/config"
	"github.com/kirban/potato-db/internal/helpers"
	loggerModule "github.com/kirban/potato-db/internal/logger"
	"github.com/kirban/potato-db/internal/network"
	"go.uber.org/zap"
	"log"
	"net"
	"os"
	"syscall"
	"time"
)

const DefaultConfigPath = "config.potato.yaml"

type AppCli struct {
	client *network.TCPClient
	logger *zap.Logger
	config *config.Config
}

func NewAppCli() (*AppCli, error) {
	app := &AppCli{}

	err := app.initDeps()

	if err != nil {
		return nil, err
	}

	return app, nil
}

func (app *AppCli) Run() error {
	defer func() {
		if err := app.logger.Sync(); err != nil {
			// ожидаемые ошибки stdout/stderr на macOS можно игнорировать
			log.Printf("logger sync error: %v", err)
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Enter command and then press enter\n")

	for {
		fmt.Printf("> ")
		query, err := reader.ReadString('\n')

		if errors.Is(err, syscall.EPIPE) {
			app.logger.Fatal("connection was closed", zap.Error(err))
		} else if err != nil {
			app.logger.Error("failed to read query", zap.Error(err))
		}

		result, err := app.client.Send([]byte(query))

		if errors.Is(err, syscall.EPIPE) {
			app.logger.Fatal("connection was closed", zap.Error(err))
			return errors.New("connection was closed")
		} else if err != nil {
			app.logger.Error("failed to send query", zap.Error(err))
			return errors.New("failed to send query")
		}

		fmt.Println(string(result))
	}
}

func (app *AppCli) initDeps() error {
	deps := []func() error{
		app.initConfig,
		app.initLogger,
		app.initClient,
	}

	for _, dep := range deps {
		if err := dep(); err != nil {
			return err
		}
	}

	return nil
}

func (app *AppCli) initConfig() error {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		configPath = DefaultConfigPath
	}

	cfg, err := config.NewConfig(configPath)

	if err != nil {
		return err
	}

	app.config = cfg
	return nil
}

func (app *AppCli) initLogger() error {
	logger, err := loggerModule.NewLogger(app.config)

	if err != nil {
		log.Fatal(err)
	}

	app.logger = logger
	return nil
}

func (app *AppCli) initClient() error {
	host := flag.String("host", "localhost", "host to connect to")
	port := flag.String("port", "8282", "port to connect to")
	idleTimeout := flag.Duration("idle-timeout", time.Minute, "idle timeout")
	maxMessageSize := flag.String("max-message-size", "4KB", "max message size")
	flag.Parse()

	maxSize, err := helpers.ParseSize(*maxMessageSize)
	if err != nil {
		app.logger.Fatal("failed to parse max message size", zap.Error(err))
		return err
	}

	addr := net.JoinHostPort(*host, *port)
	client, err := network.NewTCPClient(addr, *idleTimeout, maxSize)

	if err != nil {
		app.logger.Fatal("failed to create tcp client", zap.Error(err))
		return err
	}

	app.client = client
	return nil
}

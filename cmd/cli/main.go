package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/kirban/potato-db/internal/helpers"
	"log"
	"net"
	"os"
	"syscall"
	"time"

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

	host := flag.String("host", "localhost", "host to connect to")
	port := flag.String("port", "8282", "port to connect to")
	idleTimeout := flag.Duration("idle-timeout", time.Minute, "idle timeout")
	maxMessageSize := flag.String("max-message-size", "4KB", "max message size")
	flag.Parse()

	maxSize, err := helpers.ParseSize(*maxMessageSize)
	if err != nil {
		logger.Fatal("failed to parse max message size", zap.Error(err))
	}

	addr := net.JoinHostPort(*host, *port)
	client, err := network.NewTCPClient(addr, *idleTimeout, maxSize)

	if err != nil {
		logger.Fatal("failed to create tcp client", zap.Error(err))
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Enter command and then press enter\n")

	for {
		fmt.Printf("> ")
		query, err := reader.ReadString('\n')

		if errors.Is(err, syscall.EPIPE) {
			logger.Fatal("connection was closed", zap.Error(err))
		} else if err != nil {
			logger.Error("failed to read query", zap.Error(err))
		}

		result, err := client.Send([]byte(query))

		if errors.Is(err, syscall.EPIPE) {
			logger.Fatal("connection was closed", zap.Error(err))
		} else if err != nil {
			logger.Error("failed to send query", zap.Error(err))
		}

		fmt.Println(string(result))
	}
}

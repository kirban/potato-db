package main

import (
	"bufio"
	"fmt"
	"github.com/kirban/potato-db/internal/config"
	"github.com/kirban/potato-db/internal/db"
	loggerModule "github.com/kirban/potato-db/internal/logger"
	"log"
	"os"
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

	database := db.NewDbBuilder(logger).
		InitStorage().
		InitCompute().
		Build()

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Enter command and then press enter\n")

	for {
		fmt.Printf("> ")
		query, _ := reader.ReadString('\n')

		result := database.ExecuteQuery(query)

		fmt.Println(result)
	}
}

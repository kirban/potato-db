package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/kirban/potato-db/internal/db"
	"github.com/kirban/potato-db/internal/db/compute"
	"github.com/kirban/potato-db/internal/db/storage"
	inmemory "github.com/kirban/potato-db/internal/db/storage/engines/in-memory"
	"go.uber.org/zap"
	"log"
	"os"
	"syscall"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			log.Fatalf("can't sync zap logger: %v", err)
		}
	}(logger)

	database, err := initDb(logger)

	if err != nil {
		logger.Fatal("can't initialize database", zap.Error(err))
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Enter command and then press enter\n")

	for {
		fmt.Printf("> ")
		query, _ := reader.ReadString('\n')

		if errors.Is(err, syscall.EPIPE) {
			logger.Fatal("connection was closed", zap.Error(err))
		} else if err != nil {
			logger.Error("failed to read query", zap.Error(err))
		}

		result := database.ExecuteQuery(query)

		fmt.Println(result)
	}
}

func initStorage(logger *zap.Logger) (*storage.Storage, error) {
	engine, err := inmemory.NewInMemoryEngine(logger)

	if err != nil {
		return nil, err
	}

	st, err := storage.NewStorage(engine, logger)

	if err != nil {
		return nil, err
	}

	return st, nil
}

func initDb(logger *zap.Logger) (*db.Database, error) {
	p := compute.NewQueryParser(logger)
	c := compute.NewCompute(p)
	st, stErr := initStorage(logger)

	if stErr != nil {
		logger.Fatal("can't initialize storage", zap.Error(stErr))
		return nil, stErr
	}

	database, err := db.NewDatabase(c, st, logger)

	if err != nil {
		logger.Fatal("can't initialize db", zap.Error(err))
		return nil, err
	}

	return database, nil
}

package db

import (
	"github.com/kirban/potato-db/internal/db/compute"
	"github.com/kirban/potato-db/internal/db/storage"
	inmemory "github.com/kirban/potato-db/internal/db/storage/engines/in-memory"
	"go.uber.org/zap"
	"log"
)

type DatabaseBuilder interface {
	InitLogger() DatabaseBuilder
	InitStorage() DatabaseBuilder
	InitCompute() DatabaseBuilder
	Build() *Database
}

type dbBuilder struct {
	logger  *zap.Logger
	storage *storage.Storage
	compute *compute.Compute
}

func NewDbBuilder() DatabaseBuilder {
	return &dbBuilder{}
}

func (d *dbBuilder) InitLogger() DatabaseBuilder {
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

	d.logger = logger
	return d
}

func (d *dbBuilder) InitStorage() DatabaseBuilder {
	engine, _ := inmemory.NewInMemoryEngine()

	d.storage = storage.
		NewDatabaseStorageBuilder().
		InitEngine(engine).
		Build()

	return d
}

func (d *dbBuilder) InitCompute() DatabaseBuilder {
	var defaultParser = compute.NewQueryParser()

	d.compute = compute.
		NewDatabaseComputeBuilder().
		InitParser(defaultParser).
		Build()

	return d
}

func (d *dbBuilder) Build() *Database {
	database, err := NewDatabase(d.compute, d.storage, d.logger)

	if err != nil {
		d.logger.Error("can't initialize database", zap.Error(err))
		return nil
	}

	return database
}

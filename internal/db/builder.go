package db

import (
	"github.com/kirban/potato-db/internal/db/compute"
	"github.com/kirban/potato-db/internal/db/storage"
	inmemory "github.com/kirban/potato-db/internal/db/storage/engines/in-memory"
	"go.uber.org/zap"
)

type DatabaseBuilder interface {
	InitStorage() DatabaseBuilder
	InitCompute() DatabaseBuilder
	Build() *Database
}

type dbBuilder struct {
	logger  *zap.Logger
	storage *storage.Storage
	compute *compute.Compute
}

func NewDbBuilder(logger *zap.Logger) DatabaseBuilder {
	return &dbBuilder{
		logger: logger,
	}
}

func (d *dbBuilder) InitStorage() DatabaseBuilder {
	engine, _ := inmemory.NewInMemoryEngine(d.logger)

	d.storage = storage.
		NewDatabaseStorageBuilder(d.logger).
		InitEngine(engine).
		Build()

	return d
}

func (d *dbBuilder) InitCompute() DatabaseBuilder {
	var defaultParser = compute.NewQueryParser(d.logger)

	d.compute = compute.
		NewDatabaseComputeBuilder(d.logger).
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

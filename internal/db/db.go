package db

import "go.uber.org/zap"

type DbExecutable interface {
	ExecuteQuery(q []byte) (any, error)
}

type computeModule interface {
	Parse(q []byte) (any, error)
}

type storageModule interface {
	Set(k string, v string) error
	Get(k string) error
	Del(k string) error
}

type Database struct {
	logger        *zap.Logger
	computeModule computeModule
	storageModule storageModule
}

func New(computeModule computeModule, storageModule storageModule, logger *zap.Logger) (*Database, error) {
	if logger == nil {
		return nil, ErrLoggerNotInitialized
	}

	if computeModule == nil {
		return nil, ErrComputeModuleNotInitialized
	}

	if storageModule == nil {
		return nil, ErrStorageModuleNotInitialized
	}

	return &Database{
		logger:        logger,
		computeModule: computeModule,
		storageModule: storageModule,
	}, nil
}

func (db *Database) ExecuteQuery(q []byte) (any, error) {
	// todo: implement execute query func
	return nil, nil
}

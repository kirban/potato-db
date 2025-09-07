package db

import (
	"errors"
	"fmt"

	"github.com/kirban/potato-db/internal/db/compute"
	"go.uber.org/zap"
)

var (
	ErrLoggerNotInitialized        = errors.New("logger is not initialized")
	ErrComputeModuleNotInitialized = errors.New("compute module is not initialized")
	ErrStorageModuleNotInitialized = errors.New("storage module is not initialized")
	ErrParseFailed                 = errors.New("failed to parse command")
	ErrUnknownCommandInQuery       = errors.New("unknown command")
)

type DbExecutable interface {
	ExecuteQuery(q string) (any, error)
}

type computeModule interface {
	Compute(q string) (compute.Query, error)
}

type storageModule interface {
	Set(k string, v string) error
	Get(k string) (string, error)
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

func (db *Database) ExecuteQuery(q string) string {
	query, err := db.computeModule.Compute(q)

	if err != nil {
		return fmt.Sprintf("%s %s", compute.QueryErrorResult, err.Error())
	}

	switch query.CommandType {
	case compute.GetCommand:
		value, err := db.storageModule.Get(query.Arguments[0])
		if err != nil {
			return fmt.Sprintf("%s %s", compute.QueryErrorResult, err.Error())
		}
		return fmt.Sprintf("%s %s", compute.QueryOkResult, value)
	case compute.SetCommand:
		if err := db.storageModule.Set(query.Arguments[0], query.Arguments[1]); err != nil {
			return fmt.Sprintf("%s %s", compute.QueryErrorResult, err.Error())
		}
		return fmt.Sprint(compute.QueryOkResult)
	case compute.DelCommand:
		if err := db.storageModule.Del(query.Arguments[0]); err != nil {
			return fmt.Sprint(compute.QueryErrorResult, err.Error())
		}
		return fmt.Sprint(compute.QueryOkResult)
	}

	return fmt.Sprintf("%s", compute.QueryErrorResult)
}

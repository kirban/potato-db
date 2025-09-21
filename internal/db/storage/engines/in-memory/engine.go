package inmemory

import (
	"errors"
	"go.uber.org/zap"
)

var (
	ErrInvalidLogger = errors.New("invalid logger")
)

type InMemEngine struct {
	dataStorage *HashTable
	logger      *zap.Logger
}

func (e *InMemEngine) Get(key string) (string, bool) {
	val, exists := e.dataStorage.Get(key)

	return val, exists
}

func (e *InMemEngine) Set(key string, value string) error {
	e.dataStorage.Set(key, value)
	return nil
}

func (e *InMemEngine) Delete(key string) error {
	e.dataStorage.Del(key)
	return nil
}

func NewInMemoryEngine(logger *zap.Logger) (*InMemEngine, error) {
	if logger == nil {
		return nil, ErrInvalidLogger
	}

	return &InMemEngine{
		dataStorage: NewHashTable(),
		logger:      logger,
	}, nil
}

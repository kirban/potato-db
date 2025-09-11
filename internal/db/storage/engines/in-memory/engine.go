package inmemory

import (
	"errors"
	"go.uber.org/zap"
)

var (
	ErrInvalidLogger = errors.New("invalid logger")
)

type Engine struct {
	logger      *zap.Logger
	dataStorage *HashTable
}

func (e *Engine) Get(key string) (string, bool) {
	val, exists := e.dataStorage.Get(key)

	return val, exists
}

func (e *Engine) Set(key string, value string) error {
	e.dataStorage.Set(key, value)
	return nil
}

func (e *Engine) Delete(key string) error {
	e.dataStorage.Del(key)
	return nil
}

func NewInMemoryEngine(logger *zap.Logger) (*Engine, error) {
	if logger == nil {
		return nil, ErrInvalidLogger
	}

	return &Engine{
		logger:      logger,
		dataStorage: NewHashTable(),
	}, nil
}

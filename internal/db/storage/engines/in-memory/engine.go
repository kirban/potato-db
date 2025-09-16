package inmemory

import (
	"errors"
)

var (
	ErrInvalidLogger = errors.New("invalid logger")
)

type InMemEngine struct {
	dataStorage *HashTable
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

func NewInMemoryEngine() (*InMemEngine, error) {
	return &InMemEngine{
		dataStorage: NewHashTable(),
	}, nil
}

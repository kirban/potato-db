package inmemory

import "errors"

var (
	ErrKeyNotFound = errors.New("key not found")
)

type Engine struct {
	dataStorage *HashTable
}

func (e *Engine) Get(key string) (string, error) {
	val, exists := e.dataStorage.Get(key)

	if !exists {
		return "", ErrKeyNotFound
	}

	return val, nil
}

func (e *Engine) Set(key string, value string) error {
	e.dataStorage.Set(key, value)
	return nil
}

func (e *Engine) Delete(key string) error {
	e.dataStorage.Del(key)
	return nil
}

func NewInMemoryEngine() *Engine {
	return &Engine{
		dataStorage: NewHashTable(),
	}
}

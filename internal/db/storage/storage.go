package storage

import (
	"errors"
	"go.uber.org/zap"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

type Engine interface {
	Get(key string) (string, bool)
	Set(key string, value string) error
	Delete(key string) error
}

type Storage struct {
	engine *Engine
	logger *zap.Logger
}

func (s *Storage) Get(key string) (string, error) {
	val, exists := (*s.engine).Get(key)

	if !exists {
		return "", ErrKeyNotFound
	}

	return val, nil
}

func (s *Storage) Set(key string, value string) error {
	err := (*s.engine).Set(key, value)

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Del(key string) error {
	err := (*s.engine).Delete(key)

	if err != nil {
		return err
	}

	return nil
}

func NewStorage(engine *Engine, logger *zap.Logger) (*Storage, error) {
	if engine == nil {
		return nil, errors.New("engine is required")
	}

	return &Storage{
		engine: engine,
		logger: logger,
	}, nil
}

package storage

import "go.uber.org/zap"

type DatabaseStorageBuilder interface {
	InitEngine(engine Engine) DatabaseStorageBuilder
	Build() *Storage
}

type dbStorageBuilder struct {
	engine *Engine
	logger *zap.Logger
}

func NewDatabaseStorageBuilder(logger *zap.Logger) DatabaseStorageBuilder {
	return &dbStorageBuilder{
		logger: logger,
	}
}

func (sb *dbStorageBuilder) InitEngine(engine Engine) DatabaseStorageBuilder {
	sb.engine = &engine
	return sb
}

func (sb *dbStorageBuilder) Build() *Storage {
	s, _ := NewStorage(sb.engine, sb.logger)
	return s
}

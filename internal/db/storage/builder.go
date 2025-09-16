package storage

type DatabaseStorageBuilder interface {
	InitEngine(engine Engine) DatabaseStorageBuilder
	Build() *Storage
}

type dbStorageBuilder struct {
	engine *Engine
}

func NewDatabaseStorageBuilder() DatabaseStorageBuilder {
	return &dbStorageBuilder{}
}

func (sb *dbStorageBuilder) InitEngine(engine Engine) DatabaseStorageBuilder {
	sb.engine = &engine
	return sb
}

func (sb *dbStorageBuilder) Build() *Storage {
	s, _ := NewStorage(sb.engine)
	return s
}

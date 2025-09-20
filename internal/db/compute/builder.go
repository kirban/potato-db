package compute

import "go.uber.org/zap"

type DatabaseComputeBuilder interface {
	InitParser(parser Parser) DatabaseComputeBuilder
	Build() *Compute
}

type dbComputeBuilder struct {
	logger *zap.Logger
	parser *Parser
}

func NewDatabaseComputeBuilder(logger *zap.Logger) DatabaseComputeBuilder {
	return &dbComputeBuilder{
		logger: logger,
	}
}

func (cb *dbComputeBuilder) InitParser(parser Parser) DatabaseComputeBuilder {
	cb.parser = &parser
	return cb
}

func (cb *dbComputeBuilder) Build() *Compute {
	return NewCompute(cb.parser)
}

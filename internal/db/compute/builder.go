package compute

type DatabaseComputeBuilder interface {
	InitParser(parser Parser) DatabaseComputeBuilder
	Build() *Compute
}

type dbComputeBuilder struct {
	parser *Parser
}

func NewDatabaseComputeBuilder() DatabaseComputeBuilder {
	return &dbComputeBuilder{}
}

func (cb *dbComputeBuilder) InitParser(parser Parser) DatabaseComputeBuilder {
	cb.parser = &parser
	return cb
}

func (cb *dbComputeBuilder) Build() *Compute {
	return NewCompute(cb.parser)
}

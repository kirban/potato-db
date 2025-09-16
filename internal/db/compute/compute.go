package compute

import "fmt"

type Compute struct {
	parser *Parser
}

func (c *Compute) Compute(q string) (*Query, error) {
	if c.parser == nil {
		return nil, fmt.Errorf("parser is not initialized")
	}

	query, err := (*c.parser).Parse(q)
	if err != nil {
		return nil, err
	}

	return query, nil
}

func NewCompute(parser *Parser) *Compute {
	return &Compute{
		parser: parser,
	}
}

package compute

type Compute struct {
	parser *QueryParser
}

func (c *Compute) Compute(q string) (*Query, error) {
	query, err := c.parser.Parse(q)

	if err != nil {
		return nil, err
	}

	return query, nil
}

func NewCompute(parser *QueryParser) *Compute {
	return &Compute{
		parser: parser,
	}
}

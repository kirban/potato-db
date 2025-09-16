package compute

import "regexp"

type Query struct {
	Arguments   []string
	CommandType CommandType
}

var QueryArgsRegExp *regexp.Regexp = regexp.MustCompile(`[^\s]+`)

type QueryResult string

var (
	QueryOkResult    QueryResult = "[ok]"
	QueryErrorResult QueryResult = "[err]"
)

type CommandType string

var (
	SetCommand CommandType = "SET"
	GetCommand CommandType = "GET"
	DelCommand CommandType = "DEL"
)

func NewQuery(c CommandType, args []string) *Query {
	return &Query{
		CommandType: c,
		Arguments:   args,
	}
}

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

func IsValidCommand(q string) bool {
	switch q {
	case string(SetCommand), string(GetCommand), string(DelCommand):
		return true
	default:
		return false
	}
}

func IsValidArguments(args []string) bool {
	// todo: implement args validation
	return true
}

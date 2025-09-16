package compute

import (
	"errors"
	"strings"
)

var (
	ErrUnknownCommand = errors.New("parse error: unknown command")
	ErrWrongNOfArgs   = errors.New("parse error: invalid number of arguments")
	ErrInvalidQuery   = errors.New("parse error: invalid query")
)

type Parser interface {
	Parse(data string) (*Query, error)
}

type QueryParser struct{}

func (q *QueryParser) Parse(data string) (*Query, error) {
	trimmed := strings.TrimSpace(data)

	if len(trimmed) == 0 {
		return nil, ErrInvalidQuery
	}

	command, err := parseCommandType(trimmed)

	if err != nil {
		return nil, err
	}

	args, err := parseArguments(trimmed)

	if err != nil {
		return nil, err
	}

	return NewQuery(command, args), nil
}

func parseCommandType(q string) (CommandType, error) {
	rawCommand := QueryArgsRegExp.FindAllString(q, -1)[0]

	switch rawCommand {
	case string(GetCommand), string(SetCommand), string(DelCommand):
		return CommandType(rawCommand), nil
	default:
		return "", ErrUnknownCommand
	}
}

func parseArguments(q string) ([]string, error) {
	splittedQuery := QueryArgsRegExp.FindAllString(q, -1)
	rawCommand, rawArgs := splittedQuery[0], splittedQuery[1:]

	switch rawCommand {
	case string(GetCommand), string(DelCommand):
		if len(rawArgs) != 1 {
			return nil, ErrWrongNOfArgs
		}
	case string(SetCommand):
		if len(rawArgs) != 2 {
			return nil, ErrWrongNOfArgs
		}
	}

	return rawArgs, nil
}

func NewQueryParser() *QueryParser {
	return &QueryParser{}
}

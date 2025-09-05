package compute

import (
	"errors"
	"strings"

	"go.uber.org/zap"
)

var (
	ErrUnknownCommand = errors.New("parse error: unknown command")
	ErrWrongNOfArgs   = errors.New("parse error: invalid number of arguments")
	ErrInvalidArgs    = errors.New("parse error: invalid args passed")
)

type Parser interface {
	Parse(data string) (Query, error)
}

type QueryParser struct {
	logger *zap.Logger
}

func (q *QueryParser) Parse(data string) (Query, error) {
	var result Query

	trimmed := strings.TrimSpace(data)
	command, err := parseCommandType(trimmed)

	if err != nil {
		return Query{}, err
	}

	result.CommandType = command

	args, err := parseArguments(trimmed)

	if err != nil {
		return Query{}, err
	}

	result.Arguments = args

	return result, nil
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

	isValid := IsValidArguments(rawArgs)

	if !isValid {
		return nil, ErrInvalidArgs
	}

	return rawArgs, nil
}

func NewQueryParser(logger *zap.Logger) Parser {
	return &QueryParser{
		logger: logger,
	}
}

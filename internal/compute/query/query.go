package query

import (
	"errors"
	"slices"
	"strings"
)

type Command string

const (
	CommandGET Command = "GET"
	CommandSET Command = "SET"
	CommandDEL Command = "DEL"
)

type Query struct {
	command Command
	args    []string
}

func NewQuery(c Command, args []string) *Query {
	return &Query{
		command: c,
		args:    args,
	}
}

func (q *Query) Command() Command {
	return q.command
}

func (q *Query) Arg(i int) string {
	return q.args[i]
}

var (
	ErrUnknownCommand = errors.New("unknown command")
	ErrInvalidQuery   = errors.New("invalid query")
)

func ParseQueryStr(queryStr string) (*Query, error) {
	split := strings.Split(queryStr, " ")
	if len(split) < 2 || len(split) > 3 {
		return nil, ErrInvalidQuery
	}

	command := Command(strings.ToUpper(split[0]))
	allCommands := []Command{CommandGET, CommandSET, CommandDEL}
	if !slices.Contains(allCommands, Command(split[0])) {
		return nil, ErrUnknownCommand
	}

	if len(split) == 3 && command != CommandSET ||
		len(split) == 2 && command == CommandSET {
		return nil, ErrInvalidQuery
	}

	return NewQuery(Command(command), split[1:]), nil
}

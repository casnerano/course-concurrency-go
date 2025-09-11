package compute

import (
	"errors"
	"fmt"
	"strings"

	"github.com/casnerano/course-concurrency-go/internal/logger"
	"github.com/casnerano/course-concurrency-go/internal/types"
)

var (
	ErrEmptyQuery         = errors.New("empty query")
	ErrInvalidQuerySyntax = errors.New("invalid query syntax")
	ErrUnknownCommand     = errors.New("unknown query command")
)

const commandOptionPrefix = "--"

var requiredCommandArgsCount = map[types.Command]int{
	types.CommandGet:   1,
	types.CommandSet:   2,
	types.CommandDel:   1,
	types.CommandClear: 0,
}

type Compute struct{}

func New() *Compute {
	return &Compute{}
}

type parsedQuery struct {
	args    []string
	options types.Options
	command types.Command
}

func (c *Compute) Parse(rawQuery string) (types.Query, error) {
	logger.Debug("input query: " + rawQuery)

	query, err := c.parse(rawQuery)
	if err != nil {
		logger.Debug("failed parse input query: " + err.Error())

		return nil, fmt.Errorf("failed parse query: %w", err)
	}

	logger.Debug(fmt.Sprintf("success parse input query: %+v", query))

	switch query.command {
	case types.CommandGet:
		return c.buildQueryGet(query)
	case types.CommandSet:
		return c.buildQuerySet(query)
	case types.CommandDel:
		return c.buildQueryDel(query)
	case types.CommandClear:
		return c.buildQueryClear(query)
	default:
		return nil, ErrUnknownCommand
	}
}

func (c *Compute) parse(rawQuery string) (parsedQuery, error) {
	tokens := strings.Fields(strings.TrimSpace(rawQuery))

	if len(tokens) == 0 {
		return parsedQuery{}, ErrEmptyQuery
	}

	var (
		command = types.Command(tokens[0])
		options = make(types.Options, 0)
	)

	tokens = tokens[1:]

	for _, token := range tokens {
		if strings.HasPrefix(token, commandOptionPrefix) {
			options = append(options, types.Option(strings.TrimPrefix(token, commandOptionPrefix)))
		}
	}

	args := tokens[len(options):]

	if len(args) < requiredCommandArgsCount[command] {
		return parsedQuery{}, ErrInvalidQuerySyntax
	}

	query := parsedQuery{
		args:    args,
		options: options,
		command: command,
	}

	return query, nil
}

func (c *Compute) buildQueryGet(query parsedQuery) (types.Query, error) {
	return types.NewQueryGet(types.Key(query.args[0]), query.options...), nil
}

func (c *Compute) buildQuerySet(query parsedQuery) (types.Query, error) {
	value := types.Value{
		Type: types.ValueTypeString,
		Data: query.args[1],
	}

	return types.NewQuerySet(types.Key(query.args[0]), &value, query.options...), nil
}

func (c *Compute) buildQueryDel(query parsedQuery) (types.Query, error) {
	keys := make([]types.Key, 0, len(query.args))
	for _, token := range query.args {
		keys = append(keys, types.Key(token))
	}

	return types.NewQueryDel(keys, query.options...), nil
}

func (c *Compute) buildQueryClear(query parsedQuery) (types.Query, error) {
	return types.NewQueryClear(query.options...), nil
}

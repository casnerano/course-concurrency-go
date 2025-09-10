package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/casnerano/course-concurrency-go/internal/types"
)

var ErrUndefinedQueryType = errors.New("undefined query type")

type compute interface {
	Parse(rawQuery string) (types.Query, error)
}

type storage interface {
	Get(context.Context, types.Key) (*types.Value, error)
	Set(context.Context, types.Key, *types.Value) error
	Del(context.Context, []types.Key) error
	Clear(context.Context) error
}

type Database struct {
	compute compute
	storage storage
}

func New(compute compute, storage storage) *Database {
	return &Database{
		compute: compute,
		storage: storage,
	}
}

func (d *Database) HandleQuery(ctx context.Context, rawQuery string) (*types.Value, error) {
	query, err := d.compute.Parse(rawQuery)
	if err != nil {
		return nil, fmt.Errorf("failed compute query: %w", err)
	}

	switch query.Command() {
	case types.CommandGet:
		return d.handleCommandGet(ctx, query)
	case types.CommandSet:
		return nil, d.handleCommandSet(ctx, query)
	case types.CommandDel:
		return nil, d.handleCommandDel(ctx, query)
	case types.CommandClear:
		return nil, d.handleCommandClear(ctx, query)
	default:
		return nil, fmt.Errorf("unknown query command: %s", query.Command())
	}
}

func (d *Database) handleCommandGet(ctx context.Context, query types.Query) (*types.Value, error) {
	queryGet, ok := query.(*types.QueryGet)
	if !ok {
		return nil, ErrUndefinedQueryType
	}

	return d.storage.Get(ctx, queryGet.Key())
}

func (d *Database) handleCommandSet(ctx context.Context, query types.Query) error {
	querySet, ok := query.(*types.QuerySet)
	if !ok {
		return ErrUndefinedQueryType
	}

	return d.storage.Set(ctx, querySet.Key(), querySet.Value())
}

func (d *Database) handleCommandDel(ctx context.Context, query types.Query) error {
	queryDel, ok := query.(*types.QueryDel)
	if !ok {
		return ErrUndefinedQueryType
	}

	return d.storage.Del(ctx, queryDel.Keys())
}

func (d *Database) handleCommandClear(ctx context.Context, query types.Query) error {
	_, ok := query.(*types.QueryClear)
	if !ok {
		return ErrUndefinedQueryType
	}

	return d.storage.Clear(ctx)
}

package database

import (
	"context"
	"fmt"

	"github.com/casnerano/course-concurrency-go/internal/types"
)

type storage interface {
	Get(context.Context, types.Key) (*types.Value, error)
	Set(context.Context, types.Key, *types.Value) error
	Del(context.Context, types.Key) error
	Clear(context.Context) error
}

type Database struct {
	storage storage
}

func New(storage storage) *Database {
	return &Database{
		storage: storage,
	}
}

type Query struct {
	Command types.Command
	Key     *types.Key
	Value   *types.Value
}

func (d *Database) HandleQuery(ctx context.Context, query Query) (*types.Value, error) {
	if err := d.validateQuery(query); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	switch query.Command {
	case types.CommandGet:
		return d.handleCommandGet(ctx, *query.Key)
	case types.CommandSet:
		return nil, d.handleCommandSet(ctx, *query.Key, query.Value)
	case types.CommandDel:
		return nil, d.handleCommandDel(ctx, *query.Key)
	case types.CommandClear:
		return nil, d.handleCommandClear(ctx)
	default:
		return nil, fmt.Errorf("unknown command: %s", query.Command)
	}
}

func (d *Database) validateQuery(query Query) error {
	if query.Command.NeedKeyArg() && query.Key == nil {
		return fmt.Errorf("key is nil")
	}

	if query.Command.NeedValueArg() && query.Value == nil {
		return fmt.Errorf("value is nil")
	}

	return nil
}

func (d *Database) handleCommandGet(ctx context.Context, key types.Key) (*types.Value, error) {
	return d.storage.Get(ctx, key)
}

func (d *Database) handleCommandSet(ctx context.Context, key types.Key, value *types.Value) error {
	return d.storage.Set(ctx, key, value)
}

func (d *Database) handleCommandDel(ctx context.Context, key types.Key) error {
	return d.storage.Del(ctx, key)
}

func (d *Database) handleCommandClear(ctx context.Context) error {
	return d.storage.Clear(ctx)
}

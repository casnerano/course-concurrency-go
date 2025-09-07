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

type Request struct {
	Command types.Command
	Key     *types.Key
	Value   *types.Value
}

func (d *Database) HandleRequest(ctx context.Context, req Request) (*types.Value, error) {
	if err := d.validateRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	switch req.Command {
	case types.CommandGet:
		return d.handleCommandGet(ctx, *req.Key)
	case types.CommandSet:
		return nil, d.handleCommandSet(ctx, *req.Key, req.Value)
	case types.CommandDel:
		return nil, d.handleCommandDel(ctx, *req.Key)
	case types.CommandClear:
		return nil, d.handleCommandClear(ctx)
	default:
		return nil, fmt.Errorf("unknown command: %s", req.Command)
	}
}

func (d *Database) validateRequest(req Request) error {
	if req.Command.NeedKeyArg() && req.Key == nil {
		return fmt.Errorf("key is nil")
	}

	if req.Command.NeedValueArg() && req.Value == nil {
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

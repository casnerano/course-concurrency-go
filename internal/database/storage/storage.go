package storage

import (
	"context"

	"github.com/casnerano/course-concurrency-go/internal/types"
)

type engine interface {
	Get(context.Context, types.Key) (*types.Value, error)
	Set(context.Context, types.Key, *types.Value) error
	Del(context.Context, []types.Key) error
	Clear(context.Context) error
}

type options struct{}

type Option func(*options)

type Storage struct {
	engine engine
}

func New(engine engine, opts ...Option) *Storage {
	defOptions := options{}

	for _, opt := range opts {
		opt(&defOptions)
	}

	return &Storage{
		engine: engine,
	}
}

func (s *Storage) Get(ctx context.Context, key types.Key) (*types.Value, error) {
	return s.engine.Get(ctx, key)
}

func (s *Storage) Set(ctx context.Context, key types.Key, value *types.Value) error {
	return s.engine.Set(ctx, key, value)
}

func (s *Storage) Del(ctx context.Context, keys []types.Key) error {
	return s.engine.Del(ctx, keys)
}

func (s *Storage) Clear(ctx context.Context) error {
	return s.engine.Clear(ctx)
}

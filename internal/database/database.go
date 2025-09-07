package database

import (
	"context"

	"github.com/casnerano/course-concurrency-go/internal/types"
)

type storage interface {
	Get(context.Context, types.Key) (types.Value, error)
	Set(context.Context, types.Key) error
	Del(context.Context, types.Key) error
}

type database struct {
	storage storage
}

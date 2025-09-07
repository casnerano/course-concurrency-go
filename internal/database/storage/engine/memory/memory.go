package memory

import (
	"context"
	"sync"

	"github.com/casnerano/course-concurrency-go/internal/types"
)

type Memory struct {
	mu    sync.RWMutex
	store map[types.Key]*types.Value
}

func New() *Memory {
	return &Memory{
		store: make(map[types.Key]*types.Value),
	}
}

func (m *Memory) Get(_ context.Context, key types.Key) (*types.Value, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, ok := m.store[key]
	if ok {
		return value, nil
	}

	return nil, nil
}

func (m *Memory) Set(_ context.Context, key types.Key, value *types.Value) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.store[key] = value

	return nil
}

func (m *Memory) Del(_ context.Context, key types.Key) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.store, key)

	return nil
}

func (m *Memory) Clear(_ context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	clear(m.store)

	return nil
}

package main

import (
	"github.com/casnerano/course-concurrency-go/internal/database"
	"github.com/casnerano/course-concurrency-go/internal/database/storage"
	"github.com/casnerano/course-concurrency-go/internal/database/storage/engine/memory"
)

func getDatabase() *database.Database {
	memoryStorage := storage.New(memory.New())
	return database.New(memoryStorage)
}

package storage

import (
	"sync"
	"time"
)

type item struct {
	kind valueType

	value string

	list []string

	expiresAt *time.Time
}

type MemoryStorage struct {
	mu   sync.RWMutex
	data map[string]item

	waiters map[string][]*blpopWaiter
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data:    make(map[string]item),
		waiters: make(map[string][]*blpopWaiter),
	}
}

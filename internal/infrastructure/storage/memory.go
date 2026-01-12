package storage

import (
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/domain/value"
)

type record struct {
	v         value.Value
	expiresAt *time.Time
}

type MemoryStorage struct {
	mu   sync.RWMutex
	data map[string]record

	waiters map[string][]*blpopWaiter
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data:    make(map[string]record),
		waiters: make(map[string][]*blpopWaiter),
	}
}

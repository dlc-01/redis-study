package storage

import (
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/ports"
)

type item struct {
	value     string
	expiresAt *time.Time
}

type MemoryStorage struct {
	mu   sync.RWMutex
	data map[string]item
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string]item),
	}
}

func (s *MemoryStorage) Set(key string, value string, opts ports.SetOptions) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var expiresAt *time.Time
	if opts.PX != nil {
		t := time.Now().Add(*opts.PX)
		expiresAt = &t
	}

	s.data[key] = item{
		value:     value,
		expiresAt: expiresAt,
	}

	return nil
}

func (s *MemoryStorage) Get(key string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	it, ok := s.data[key]
	if !ok {
		return "", false, nil
	}

	if it.expiresAt != nil && time.Now().After(*it.expiresAt) {
		delete(s.data, key)
		return "", false, nil
	}

	return it.value, true, nil
}

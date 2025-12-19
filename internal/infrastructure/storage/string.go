package storage

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/ports"
)

func (s *MemoryStorage) Set(key, value string, opts ports.SetOptions) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var expiresAt *time.Time
	if opts.PX != nil {
		t := time.Now().Add(*opts.PX)
		expiresAt = &t
	}

	s.data[key] = item{
		kind:      typeString,
		value:     value,
		expiresAt: expiresAt,
	}
	return nil
}

func (s *MemoryStorage) Get(key string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	it, ok := s.data[key]
	if !ok || s.cleanupIfExpired(key, it) {
		return "", false, nil
	}

	if it.kind != typeString {
		return "", false, ErrWrongType
	}

	return it.value, true, nil
}

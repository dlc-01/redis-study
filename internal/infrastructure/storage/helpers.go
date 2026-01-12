package storage

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/domain/value"
)

func (s *MemoryStorage) isExpired(r record) bool {
	return r.expiresAt != nil && time.Now().After(*r.expiresAt)
}

func (s *MemoryStorage) cleanupIfExpired(key string, r record) bool {
	if s.isExpired(r) {
		delete(s.data, key)
		return true
	}
	return false
}
func (s *MemoryStorage) getRecordLocked(key string) (record, bool) {
	r, ok := s.data[key]
	if !ok {
		return record{}, false
	}

	if r.expiresAt != nil && time.Now().After(*r.expiresAt) {
		delete(s.data, key)
		return record{}, false
	}

	return r, true
}

func (s *MemoryStorage) setRecordLocked(key string, v value.Value, expiresAt *time.Time) {
	s.data[key] = record{v: v, expiresAt: expiresAt}
}

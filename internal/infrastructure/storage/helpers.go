package storage

import (
	"errors"
	"time"
)

var ErrWrongType = errors.New(
	"WRONGTYPE Operation against a key holding the wrong kind of value",
)

func (s *MemoryStorage) isExpired(it item) bool {
	return it.expiresAt != nil && time.Now().After(*it.expiresAt)
}

func (s *MemoryStorage) cleanupIfExpired(key string, it item) bool {
	if s.isExpired(it) {
		delete(s.data, key)
		return true
	}
	return false
}

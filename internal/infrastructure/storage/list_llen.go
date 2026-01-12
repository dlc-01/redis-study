package storage

import (
	"github.com/codecrafters-io/redis-starter-go/internal/domain/rerrors"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/value"
)

func (s *MemoryStorage) LLen(key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	r, ok := s.getRecordLocked(key)
	if !ok {
		return 0, nil
	}

	lst, ok := r.v.(value.List)
	if !ok {
		return 0, rerrors.ErrWrongType
	}

	return len(lst.V), nil
}

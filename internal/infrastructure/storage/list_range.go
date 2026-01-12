package storage

import (
	"github.com/codecrafters-io/redis-starter-go/internal/domain/rerrors"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/value"
)

func (s *MemoryStorage) LRange(key string, start, stop int) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	r, ok := s.getRecordLocked(key)
	if !ok {
		return []string{}, nil
	}

	lst, ok := r.v.(value.List)
	if !ok {
		return []string{}, rerrors.ErrWrongType
	}

	n := len(lst.V)
	if n == 0 {
		return []string{}, nil
	}

	if start < 0 {
		start = n + start
	}
	if stop < 0 {
		stop = n + stop
	}

	if start < 0 {
		start = 0
	}
	if stop >= n {
		stop = n - 1
	}

	if start > stop || start >= n {
		return []string{}, nil
	}

	result := lst.V[start : stop+1]

	out := make([]string, len(result))
	copy(out, result)

	return out, nil
}

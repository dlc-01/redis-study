package storage

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/domain/rerrors"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/value"
	"github.com/codecrafters-io/redis-starter-go/internal/ports"
)

func (s *MemoryStorage) Set(key string, v string, opts ports.SetOptions) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var exp *time.Time
	if opts.PX != nil {
		t := time.Now().Add(*opts.PX)
		exp = &t
	}

	s.setRecordLocked(key, value.String{V: v}, exp)
	return nil
}

func (s *MemoryStorage) Get(key string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	r, ok := s.getRecordLocked(key)
	if !ok {
		return "", false, nil
	}

	str, ok := r.v.(value.String)
	if !ok {
		return "", false, rerrors.ErrWrongType
	}

	return str.V, true, nil
}

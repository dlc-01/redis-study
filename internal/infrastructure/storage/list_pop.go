package storage

import (
	"github.com/codecrafters-io/redis-starter-go/internal/domain/rerrors"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/value"
)

func (s *MemoryStorage) LPop(key string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	r, ok := s.getRecordLocked(key)
	if !ok {
		return "", false, nil
	}

	lst, ok := r.v.(value.List)
	if !ok {
		return "", false, rerrors.ErrWrongType
	}

	if len(lst.V) == 0 {
		delete(s.data, key)
		return "", false, nil
	}

	v := lst.V[0]
	lst.V = lst.V[1:]

	if len(lst.V) == 0 {
		delete(s.data, key)
	} else {
		s.setRecordLocked(key, lst, r.expiresAt)
	}

	return v, true, nil
}

func (s *MemoryStorage) LPopN(key string, count int) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if count <= 0 {
		return []string{}, nil
	}

	r, ok := s.getRecordLocked(key)
	if !ok {
		return []string{}, nil
	}

	lst, ok := r.v.(value.List)
	if !ok {
		return []string{}, rerrors.ErrWrongType
	}

	if len(lst.V) == 0 {
		delete(s.data, key)
		return []string{}, nil
	}

	if count >= len(lst.V) {
		out := append([]string{}, lst.V...)
		delete(s.data, key)
		return out, nil
	}

	out := append([]string{}, lst.V[:count]...)
	lst.V = lst.V[count:]
	s.setRecordLocked(key, lst, r.expiresAt)

	return out, nil
}

func (s *MemoryStorage) RPop(key string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	r, ok := s.getRecordLocked(key)
	if !ok {
		return "", false, nil
	}

	lst, ok := r.v.(value.List)
	if !ok {
		return "", false, rerrors.ErrWrongType
	}

	if len(lst.V) == 0 {
		delete(s.data, key)
		return "", false, nil
	}

	v := lst.V[len(lst.V)-1]
	lst.V = lst.V[:len(lst.V)-1]

	if len(lst.V) == 0 {
		delete(s.data, key)
	} else {
		s.setRecordLocked(key, lst, r.expiresAt)
	}

	return v, true, nil
}

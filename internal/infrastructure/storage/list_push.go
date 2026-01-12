package storage

import (
	"github.com/codecrafters-io/redis-starter-go/internal/domain/rerrors"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/value"
)

func (s *MemoryStorage) RPush(key string, values ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	r, ok := s.getRecordLocked(key)
	var lst value.List

	if !ok {
		lst = value.List{V: []string{}}
		r = record{v: lst, expiresAt: nil}
	} else {
		var isList bool
		lst, isList = r.v.(value.List)
		if !isList {
			return 0, rerrors.ErrWrongType
		}
	}

	if ws := s.waiters[key]; len(ws) > 0 {
		waiter := ws[0]
		s.waiters[key] = ws[1:]

		val := values[0]
		waiter.ch <- val

		if len(values) > 1 {
			lst.V = append(lst.V, values[1:]...)
			s.setRecordLocked(key, lst, r.expiresAt)
		}

		return 1, nil
	}

	lst.V = append(lst.V, values...)
	s.setRecordLocked(key, lst, r.expiresAt)
	return len(lst.V), nil
}

func (s *MemoryStorage) LPush(key string, values ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	r, ok := s.getRecordLocked(key)
	var lst value.List

	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}

	if !ok {
		lst = value.List{V: []string{}}
		r = record{v: lst, expiresAt: nil}
	} else {
		var isList bool
		lst, isList = r.v.(value.List)
		if !isList {
			return 0, rerrors.ErrWrongType
		}
	}

	if ws := s.waiters[key]; len(ws) > 0 {
		waiter := ws[0]
		s.waiters[key] = ws[1:]

		val := values[0]
		waiter.ch <- val

		if len(values) > 1 {
			lst.V = append(values[1:], lst.V...)
			s.setRecordLocked(key, lst, r.expiresAt)
		}

		return 1, nil
	}

	lst.V = append(values, lst.V...)
	s.setRecordLocked(key, lst, r.expiresAt)
	return len(lst.V), nil
}

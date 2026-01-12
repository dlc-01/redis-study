package storage

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/domain/rerrors"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/value"
)

func (s *MemoryStorage) BLPop(key string) (string, bool, error) {
	s.mu.Lock()

	r, ok := s.getRecordLocked(key)
	if ok {
		lst, isList := r.v.(value.List)
		if !isList {
			s.mu.Unlock()
			return "", false, rerrors.ErrWrongType
		}

		if len(lst.V) > 0 {
			v := lst.V[0]
			lst.V = lst.V[1:]
			s.setRecordLocked(key, lst, r.expiresAt)
			s.mu.Unlock()
			return v, true, nil
		}
	}

	waiter := &blpopWaiter{ch: make(chan string, 1)}
	s.waiters[key] = append(s.waiters[key], waiter)
	s.mu.Unlock()

	v := <-waiter.ch
	return v, true, nil
}

func (s *MemoryStorage) BLPopWithTimeout(
	key string,
	timeout time.Duration,
) (string, bool, error) {

	s.mu.Lock()

	r, ok := s.getRecordLocked(key)
	if ok {
		lst, isList := r.v.(value.List)
		if !isList {
			s.mu.Unlock()
			return "", false, rerrors.ErrWrongType
		}

		if len(lst.V) > 0 {
			v := lst.V[0]
			lst.V = lst.V[1:]
			s.setRecordLocked(key, lst, r.expiresAt)
			s.mu.Unlock()
			return v, true, nil
		}
	}

	waiter := &blpopWaiter{ch: make(chan string, 1)}
	s.waiters[key] = append(s.waiters[key], waiter)
	s.mu.Unlock()

	select {
	case v := <-waiter.ch:
		return v, true, nil
	case <-time.After(timeout):
		s.mu.Lock()
		ws := s.waiters[key]
		for i, w := range ws {
			if w == waiter {
				s.waiters[key] = append(ws[:i], ws[i+1:]...)
				break
			}
		}
		s.mu.Unlock()
		return "", false, nil
	}
}

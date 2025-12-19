package storage

import "time"

func (s *MemoryStorage) BLPop(key string) (string, bool, error) {
	s.mu.Lock()

	it, ok := s.data[key]
	if ok && it.kind == typeList && len(it.list) > 0 {
		v := it.list[0]
		it.list = it.list[1:]
		s.data[key] = it
		s.mu.Unlock()
		return v, true, nil
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

	it, ok := s.data[key]
	if ok && it.kind == typeList && len(it.list) > 0 {
		v := it.list[0]
		it.list = it.list[1:]
		s.data[key] = it
		s.mu.Unlock()
		return v, true, nil
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

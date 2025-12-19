package storage

import "time"

func (s *MemoryStorage) LLen(key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	it, ok := s.data[key]
	if !ok {
		return 0, nil
	}

	if it.expiresAt != nil && time.Now().After(*it.expiresAt) {
		delete(s.data, key)
		return 0, nil
	}

	if it.kind != typeList {
		return 0, ErrWrongType
	}

	return len(it.list), nil
}

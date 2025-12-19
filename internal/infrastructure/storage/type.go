package storage

import "time"

func (s *MemoryStorage) Type(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	it, ok := s.data[key]
	if !ok {
		return "none", nil
	}

	if it.expiresAt != nil && time.Now().After(*it.expiresAt) {
		delete(s.data, key)
		return "none", nil
	}

	switch it.kind {
	case typeString:
		return "string", nil
	case typeList:
		return "list", nil
	default:
		return "none", nil
	}
}

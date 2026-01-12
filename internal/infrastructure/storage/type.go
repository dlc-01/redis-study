package storage

func (s *MemoryStorage) Type(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	r, ok := s.getRecordLocked(key)
	if !ok {
		return "none", nil
	}

	return string(r.v.Type()), nil
}

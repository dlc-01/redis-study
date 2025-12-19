package storage

func (s *MemoryStorage) LPop(key string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	it, ok := s.data[key]
	if !ok || it.kind != typeList || len(it.list) == 0 {
		return "", false, nil
	}

	v := it.list[0]
	it.list = it.list[1:]

	if len(it.list) == 0 {
		delete(s.data, key)
	} else {
		s.data[key] = it
	}

	return v, true, nil
}

func (s *MemoryStorage) LPopN(key string, count int) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if count <= 0 {
		return []string{}, nil
	}

	it, ok := s.data[key]
	if !ok || it.kind != typeList || len(it.list) == 0 {
		return []string{}, nil
	}

	if count >= len(it.list) {
		out := append([]string{}, it.list...)
		delete(s.data, key)
		return out, nil
	}

	out := append([]string{}, it.list[:count]...)
	it.list = it.list[count:]
	s.data[key] = it

	return out, nil
}

func (s *MemoryStorage) RPop(key string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	it, ok := s.data[key]
	if !ok || it.kind != typeList || len(it.list) == 0 {
		return "", false, nil
	}

	v := it.list[len(it.list)-1]
	it.list = it.list[:len(it.list)-1]

	if len(it.list) == 0 {
		delete(s.data, key)
	} else {
		s.data[key] = it
	}

	return v, true, nil
}

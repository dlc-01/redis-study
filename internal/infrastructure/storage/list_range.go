package storage

func (s *MemoryStorage) LRange(key string, start, stop int) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	it, ok := s.data[key]
	if !ok || it.kind != typeList {
		return []string{}, nil
	}

	n := len(it.list)
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

	result := it.list[start : stop+1]

	out := make([]string, len(result))
	copy(out, result)

	return out, nil
}

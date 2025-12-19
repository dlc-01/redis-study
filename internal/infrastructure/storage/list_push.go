package storage

func (s *MemoryStorage) RPush(key string, values ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	it, ok := s.data[key]
	if ok && s.cleanupIfExpired(key, it) {
		ok = false
	}

	if !ok {
		it = item{
			kind: typeList,
			list: []string{},
		}
	}

	if it.kind != typeList {
		return 0, ErrWrongType
	}

	if ws := s.waiters[key]; len(ws) > 0 {
		waiter := ws[0]
		s.waiters[key] = ws[1:]

		val := values[0]
		waiter.ch <- val

		if len(values) > 1 {
			it.list = append(it.list, values[1:]...)
			s.data[key] = it
		}

		return 1, nil
	}

	it.list = append(it.list, values...)
	s.data[key] = it
	return len(it.list), nil
}

func (s *MemoryStorage) LPush(key string, values ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	it, ok := s.data[key]
	if ok && s.cleanupIfExpired(key, it) {
		ok = false
	}

	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}

	if !ok {
		it = item{
			kind: typeList,
			list: []string{},
		}
	}

	if it.kind != typeList {
		return 0, ErrWrongType
	}

	if ws := s.waiters[key]; len(ws) > 0 {
		waiter := ws[0]
		s.waiters[key] = ws[1:]

		val := values[0]
		waiter.ch <- val

		if len(values) > 1 {
			it.list = append(values[1:], it.list...)
			s.data[key] = it
		}

		return 1, nil
	}

	it.list = append(values, it.list...)
	s.data[key] = it
	return len(it.list), nil
}

package storage

type xreadWaiter struct {
	ch chan struct{}
}

func (s *MemoryStorage) notifyStreamWaitersLocked(key string) {
	ws := s.streamWaiters[key]
	if len(ws) == 0 {
		return
	}
	delete(s.streamWaiters, key)

	for _, w := range ws {
		select {
		case w.ch <- struct{}{}:
		default:
		}
	}
}

func (s *MemoryStorage) addStreamWaiterLocked(key string, w *xreadWaiter) {
	s.streamWaiters[key] = append(s.streamWaiters[key], w)
}

func (s *MemoryStorage) removeStreamWaiterLocked(key string, w *xreadWaiter) {
	ws := s.streamWaiters[key]
	for i := 0; i < len(ws); i++ {
		if ws[i] == w {
			ws = append(ws[:i], ws[i+1:]...)
			break
		}
	}
	if len(ws) == 0 {
		delete(s.streamWaiters, key)
	} else {
		s.streamWaiters[key] = ws
	}
}

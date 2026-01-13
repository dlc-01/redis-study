package storage

import (
	"github.com/codecrafters-io/redis-starter-go/internal/application/streamidcodec"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/rerrors"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/streamid"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/value"
)

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

func (s *MemoryStorage) computeXReadStartsLocked(keys []string, ids []string) ([]streamid.ID, error) {
	starts := make([]streamid.ID, len(keys))
	for i := range keys {
		key := keys[i]
		id := ids[i]

		if id == "$" {
			last, _, err := s.lastIDLocked(key)
			if err != nil {
				return nil, err
			}
			starts[i] = last
			continue
		}

		start, err := streamidcodec.ParseRangeID(id, true)
		if err != nil {
			return nil, err
		}
		starts[i] = start
	}
	return starts, nil
}

func (s *MemoryStorage) xreadManyFromStartsLocked(keys []string, starts []streamid.ID) (map[string][]value.StreamEntry, error) {
	out := make(map[string][]value.StreamEntry)
	for i := range keys {
		entries, err := s.xreadOneLocked(keys[i], starts[i])
		if err != nil {
			return nil, err
		}
		if len(entries) > 0 {
			out[keys[i]] = entries
		}
	}
	return out, nil
}

func (s *MemoryStorage) xreadOneLocked(key string, start streamid.ID) ([]value.StreamEntry, error) {
	r, ok := s.getRecordLocked(key)
	if !ok {
		return []value.StreamEntry{}, nil
	}

	st, ok := r.v.(value.Stream)
	if !ok {
		return nil, rerrors.ErrWrongType
	}

	if len(st.V) == 0 {
		return []value.StreamEntry{}, nil
	}

	i := lowerBoundGT(st.V, start)
	if i >= len(st.V) {
		return []value.StreamEntry{}, nil
	}

	out := make([]value.StreamEntry, len(st.V)-i)
	copy(out, st.V[i:])
	return out, nil
}

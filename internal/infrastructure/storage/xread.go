package storage

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/application/streamidcodec"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/rerrors"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/streamid"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/value"
)

func (s *MemoryStorage) lastIDLocked(key string) (streamid.ID, bool, error) {
	r, ok := s.getRecordLocked(key)
	if !ok {
		return streamid.ID{Time: 0, Seq: 0}, false, nil
	}
	st, ok := r.v.(value.Stream)
	if !ok {
		return streamid.ID{}, false, rerrors.ErrWrongType
	}
	if len(st.V) == 0 {
		return streamid.ID{Time: 0, Seq: 0}, true, nil
	}

	spec, err := streamidcodec.ParseSpec(st.V[len(st.V)-1].ID)
	if err != nil {
		return streamid.ID{}, false, err
	}
	if spec.Kind != streamidcodec.SpecExplicit {
		return streamid.ID{}, false, rerrors.ErrInvalidStreamID
	}
	return streamid.ID{Time: spec.Time, Seq: spec.Seq}, true, nil
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

	out := make([]value.StreamEntry, 0)
	for _, e := range st.V {
		es, err := streamidcodec.ParseSpec(e.ID)
		if err != nil {
			return nil, err
		}
		if es.Kind != streamidcodec.SpecExplicit {
			return nil, rerrors.ErrInvalidStreamID
		}
		eid := streamid.ID{Time: es.Time, Seq: es.Seq}
		if streamid.Compare(eid, start) <= 0 {
			continue
		}
		out = append(out, e)
	}
	return out, nil
}

func (s *MemoryStorage) XReadMany(keys []string, ids []string) (map[string][]value.StreamEntry, error) {
	if len(keys) != len(ids) {
		return nil, rerrors.ErrInvalidArgs
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	out := make(map[string][]value.StreamEntry, len(keys))

	for i := range keys {
		key := keys[i]
		id := ids[i]

		var start streamid.ID
		if id == "$" {
			last, _, err := s.lastIDLocked(key)
			if err != nil {
				return nil, err
			}
			start = last
		} else {
			parsed, err := streamidcodec.ParseRangeID(id, true)
			if err != nil {
				return nil, err
			}
			start = parsed
		}

		entries, err := s.xreadOneLocked(key, start)
		if err != nil {
			return nil, err
		}
		if len(entries) > 0 {
			out[key] = entries
		}
	}

	return out, nil
}

func (s *MemoryStorage) XReadManyBlocked(keys []string, ids []string, timeout time.Duration) (map[string][]value.StreamEntry, error) {
	if len(keys) != len(ids) {
		return nil, rerrors.ErrInvalidArgs
	}

	s.mu.Lock()
	starts, err := s.computeXReadStartsLocked(keys, ids)
	if err != nil {
		s.mu.Unlock()
		return nil, err
	}

	res, err := s.xreadManyFromStartsLocked(keys, starts)
	if err != nil {
		s.mu.Unlock()
		return nil, err
	}
	if len(res) > 0 {
		s.mu.Unlock()
		return res, nil
	}

	w := &xreadWaiter{ch: make(chan struct{}, 1)}
	for _, key := range keys {
		s.addStreamWaiterLocked(key, w)
	}
	s.mu.Unlock()

	if timeout == 0 {
		<-w.ch
	} else {
		select {
		case <-w.ch:
		case <-time.After(timeout):
			s.mu.Lock()
			for _, key := range keys {
				s.removeStreamWaiterLocked(key, w)
			}
			s.mu.Unlock()
			return nil, nil
		}
	}

	s.mu.Lock()
	for _, key := range keys {
		s.removeStreamWaiterLocked(key, w)
	}

	res, err = s.xreadManyFromStartsLocked(keys, starts)
	s.mu.Unlock()

	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return map[string][]value.StreamEntry{}, nil
	}
	return res, nil
}

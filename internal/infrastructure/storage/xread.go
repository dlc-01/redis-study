package storage

import (
	"github.com/codecrafters-io/redis-starter-go/internal/application/streamidcodec"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/rerrors"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/streamid"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/value"
)

func (s *MemoryStorage) xreadOneLocked(key string, startID streamid.ID) ([]value.StreamEntry, error) {
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

		if streamid.Compare(eid, startID) <= 0 {
			continue
		}
		out = append(out, e)
	}

	return out, nil
}

func (s *MemoryStorage) XRead(key string, id string) ([]value.StreamEntry, error) {
	startID, err := streamidcodec.ParseRangeID(id, true)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	return s.xreadOneLocked(key, startID)
}

func (s *MemoryStorage) XReadMany(keys []string, ids []string) (map[string][]value.StreamEntry, error) {
	if len(keys) != len(ids) {
		return nil, rerrors.ErrInvalidArgs
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	out := make(map[string][]value.StreamEntry, len(keys))

	for i := 0; i < len(keys); i++ {
		startID, err := streamidcodec.ParseRangeID(ids[i], true)
		if err != nil {
			return nil, err
		}

		entries, err := s.xreadOneLocked(keys[i], startID)
		if err != nil {
			return nil, err
		}

		if len(entries) > 0 {
			out[keys[i]] = entries
		}
	}

	return out, nil
}

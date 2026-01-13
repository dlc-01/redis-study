package storage

import (
	"github.com/codecrafters-io/redis-starter-go/internal/application/streamidcodec"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/rerrors"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/streamid"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/value"
)

func (s *MemoryStorage) XRange(key string, start string, end string) ([]value.StreamEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	r, ok := s.getRecordLocked(key)
	if !ok {
		return []value.StreamEntry{}, nil
	}

	st, ok := r.v.(value.Stream)
	if !ok {
		return nil, rerrors.ErrWrongType
	}

	startID, err := streamidcodec.ParseRangeID(start, true)
	if err != nil {
		return nil, err
	}

	endID, err := streamidcodec.ParseRangeID(end, false)
	if err != nil {
		return nil, err
	}

	if len(st.V) == 0 {
		return []value.StreamEntry{}, nil
	}

	i := lowerBoundGE(st.V, startID)

	out := make([]value.StreamEntry, 0)
	for i < len(st.V) {
		e := st.V[i]
		if streamid.Compare(e.ID, endID) > 0 {
			break
		}
		out = append(out, e)
		i++
	}

	return out, nil
}

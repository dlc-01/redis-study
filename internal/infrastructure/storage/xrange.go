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

	out := make([]value.StreamEntry, 0)

	for _, e := range st.V {
		spec, err := streamidcodec.ParseSpec(e.ID)
		if err != nil {
			return nil, err
		}
		if spec.Kind != streamidcodec.SpecExplicit {
			return nil, rerrors.ErrInvalidStreamID
		}

		eid := streamid.ID{Time: spec.Time, Seq: spec.Seq}

		if streamid.Compare(eid, startID) < 0 {
			continue
		}
		if streamid.Compare(eid, endID) > 0 {
			break
		}

		out = append(out, e)
	}

	return out, nil
}

package storage

import (
	"github.com/codecrafters-io/redis-starter-go/internal/domain/rerrors"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/value"
)

func (s *MemoryStorage) XRange(key string, start string, end string) ([]value.StreamEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	r, ok := s.getRecordLocked(key)
	if !ok {
		return []value.StreamEntry{}, nil
	}

	stream, ok := r.v.(value.Stream)
	if !ok {
		return nil, rerrors.ErrWrongType
	}

	startID, err := parseRangeID(start, true)
	if err != nil {
		return nil, err
	}

	endID, err := parseRangeID(end, false)
	if err != nil {
		return nil, err
	}

	out := make([]value.StreamEntry, 0)

	for _, e := range stream.V {
		espec, err := parseIDSpec(e.ID)
		if err != nil {
			return nil, err
		}
		eid := streamID{Time: espec.t, Seq: espec.seq}

		if cmpID(eid, startID) < 0 {
			continue
		}

		if cmpID(eid, endID) > 0 {
			break
		}

		out = append(out, e)
	}

	return out, nil
}

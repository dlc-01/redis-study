package storage

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/domain/rerrors"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/value"
)

func (s *MemoryStorage) XAdd(key string, id string, fields map[string]string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	spec, err := parseIDSpec(id)
	if err != nil {
		return "", err
	}

	r, ok := s.getRecordLocked(key)

	var (
		stream       value.Stream
		expiresAtPtr *time.Time
		lastID       *streamID
	)

	if ok {
		expiresAtPtr = r.expiresAt

		curr, isStream := r.v.(value.Stream)
		if !isStream {
			return "", rerrors.ErrWrongType
		}
		stream = curr

		if len(stream.V) > 0 {
			parsedLast, _ := parseIDSpec(stream.V[len(stream.V)-1].ID)
			tmp := streamID{Time: parsedLast.t, Seq: parsedLast.seq}
			lastID = &tmp
		}
	} else {
		stream = value.Stream{V: []value.StreamEntry{}}
	}

	var newID streamID

	switch spec.kind {
	case specExplicit:
		newID = streamID{Time: spec.t, Seq: spec.seq}

	case specAutoSeq:
		t := spec.t

		if lastID != nil && t < lastID.Time {
			return "", rerrors.ErrXAddIDTooSmall
		}

		if lastID != nil && t == lastID.Time {
			newID = streamID{Time: t, Seq: lastID.Seq + 1}
		} else {
			start := int64(0)
			if t == 0 {
				start = 1
			}
			newID = streamID{Time: t, Seq: start}
		}
	}

	if newID.Time == 0 && newID.Seq == 0 {
		return "", rerrors.ErrXAddZeroID
	}

	if lastID != nil {
		if newID.Time < lastID.Time ||
			(newID.Time == lastID.Time && newID.Seq <= lastID.Seq) {
			return "", rerrors.ErrXAddIDTooSmall
		}
	}

	idStr := newID.String()

	stream.V = append(stream.V, value.StreamEntry{
		ID:     idStr,
		Fields: fields,
	})

	s.setRecordLocked(key, stream, expiresAtPtr)
	return idStr, nil
}

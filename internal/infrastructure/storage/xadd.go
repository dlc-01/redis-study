package storage

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/domain/rerrors"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/value"
)

func (s *MemoryStorage) XAdd(
	key string,
	id string,
	fields map[string]string,
) (string, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	newID, err := parseStreamID(id)
	if err != nil {
		return "", err
	}

	if newID.Time == 0 && newID.Seq == 0 {
		return "", rerrors.ErrXAddZeroID
	}

	r, ok := s.getRecordLocked(key)

	var stream value.Stream
	var expiresAt = (*time.Time)(nil)

	if ok {
		expiresAt = r.expiresAt

		curr, isStream := r.v.(value.Stream)
		if !isStream {
			return "", rerrors.ErrWrongType
		}
		stream = curr

		if len(stream.V) > 0 {
			last := stream.V[len(stream.V)-1]
			lastID, _ := parseStreamID(last.ID)

			if newID.Time < lastID.Time ||
				(newID.Time == lastID.Time && newID.Seq <= lastID.Seq) {
				return "", rerrors.ErrXAddIDTooSmall
			}
		}
	} else {
		stream = value.Stream{V: []value.StreamEntry{}}
	}

	stream.V = append(stream.V, value.StreamEntry{
		ID:     id,
		Fields: fields,
	})

	s.setRecordLocked(key, stream, expiresAt)
	return id, nil
}

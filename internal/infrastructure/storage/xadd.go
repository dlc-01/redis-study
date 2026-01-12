package storage

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/domain/rerrors"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/value"
)

func (s *MemoryStorage) XAdd(key string, id string, fields map[string]string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	r, ok := s.getRecordLocked(key)

	var (
		stream       value.Stream
		expiresAtPtr = (*time.Time)(nil)
	)

	if ok {
		expiresAtPtr = r.expiresAt

		curr, isStream := r.v.(value.Stream)
		if !isStream {
			return "", rerrors.ErrWrongType
		}
		stream = curr
	} else {
		stream = value.Stream{V: []value.StreamEntry{}}
	}

	stream.V = append(stream.V, value.StreamEntry{
		ID:     id,
		Fields: fields,
	})

	s.setRecordLocked(key, stream, expiresAtPtr)
	return id, nil
}

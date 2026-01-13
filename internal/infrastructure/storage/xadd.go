package storage

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/application/streamidcodec"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/rerrors"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/streamid"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/value"
)

func (s *MemoryStorage) XAdd(key string, id string, values []value.StreamKV) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	spec, err := streamidcodec.ParseSpec(id)
	if err != nil {
		return "", err
	}

	r, ok := s.getRecordLocked(key)

	var (
		st           value.Stream
		expiresAtPtr *time.Time
		last         streamid.ID
		hasLast      bool
	)

	if ok {
		expiresAtPtr = r.expiresAt

		curr, isStream := r.v.(value.Stream)
		if !isStream {
			return "", rerrors.ErrWrongType
		}
		st = curr

		hasLast = len(st.V) > 0
		last = st.LastID
	} else {
		st = value.Stream{
			V:      []value.StreamEntry{},
			LastID: streamid.ID{Time: 0, Seq: 0},
		}
	}

	var newID streamid.ID

	switch spec.Kind {
	case streamidcodec.SpecExplicit:
		newID = streamid.ID{Time: spec.Time, Seq: spec.Seq}

	case streamidcodec.SpecAutoSeq:
		t := spec.Time

		if hasLast && t < last.Time {
			return "", rerrors.ErrXAddIDTooSmall
		}

		if hasLast && t == last.Time {
			newID = streamid.ID{Time: t, Seq: last.Seq + 1}
		} else {
			start := int64(0)
			if t == 0 {
				start = 1
			}
			newID = streamid.ID{Time: t, Seq: start}
		}

	case streamidcodec.SpecAutoID:
		t := time.Now().UnixMilli()
		seq := int64(0)
		if hasLast && last.Time == t {
			seq = last.Seq + 1
		}
		newID = streamid.ID{Time: t, Seq: seq}

	default:
		return "", rerrors.ErrInvalidArgs
	}

	if newID.Time == 0 && newID.Seq == 0 {
		return "", rerrors.ErrXAddZeroID
	}

	if hasLast {
		if newID.Time < last.Time || (newID.Time == last.Time && newID.Seq <= last.Seq) {
			return "", rerrors.ErrXAddIDTooSmall
		}
	}

	idStr := newID.String()

	st.V = append(st.V, value.StreamEntry{
		IDRaw:  idStr,
		ID:     newID,
		Values: values,
	})
	st.LastID = newID

	s.setRecordLocked(key, st, expiresAtPtr)
	s.notifyStreamWaitersLocked(key)

	return idStr, nil
}

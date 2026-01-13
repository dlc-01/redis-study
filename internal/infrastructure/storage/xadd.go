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
		expiresAtPtr = (*time.Time)(nil)
		last         = (*streamid.ID)(nil)
	)

	if ok {
		expiresAtPtr = r.expiresAt

		curr, isStream := r.v.(value.Stream)
		if !isStream {
			return "", rerrors.ErrWrongType
		}
		st = curr

		if n := len(st.V); n > 0 {
			lastEntryID := st.V[n-1].ID
			lastSpec, err := streamidcodec.ParseSpec(lastEntryID)
			if err != nil {
				return "", err
			}
			if lastSpec.Kind != streamidcodec.SpecExplicit {
				return "", rerrors.ErrInvalidStreamID
			}

			tmp := streamid.ID{Time: lastSpec.Time, Seq: lastSpec.Seq}
			last = &tmp
		}
	} else {
		st = value.Stream{V: []value.StreamEntry{}}
	}

	var newID streamid.ID

	switch spec.Kind {
	case streamidcodec.SpecExplicit:
		newID = streamid.ID{Time: spec.Time, Seq: spec.Seq}

	case streamidcodec.SpecAutoSeq:
		t := spec.Time

		if last != nil && t < last.Time {
			return "", rerrors.ErrXAddIDTooSmall
		}

		if last != nil && t == last.Time {
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
		if last != nil && last.Time == t {
			seq = last.Seq + 1
		}

		newID = streamid.ID{Time: t, Seq: seq}

	default:
		return "", rerrors.ErrInvalidArgs
	}

	if newID.Time == 0 && newID.Seq == 0 {
		return "", rerrors.ErrXAddZeroID
	}

	if last != nil {
		if newID.Time < last.Time || (newID.Time == last.Time && newID.Seq <= last.Seq) {
			return "", rerrors.ErrXAddIDTooSmall
		}
	}

	idStr := newID.String()

	st.V = append(st.V, value.StreamEntry{
		ID:     idStr,
		Values: values,
	})

	s.setRecordLocked(key, st, expiresAtPtr)

	s.notifyStreamWaitersLocked(key)
	return idStr, nil
}

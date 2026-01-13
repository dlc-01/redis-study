package storage

import (
	"time"

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
	return st.LastID, true, nil
}

func (s *MemoryStorage) XReadMany(keys []string, ids []string) (map[string][]value.StreamEntry, error) {
	if len(keys) != len(ids) {
		return nil, rerrors.ErrInvalidArgs
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	starts, err := s.computeXReadStartsLocked(keys, ids)
	if err != nil {
		return nil, err
	}

	return s.xreadManyFromStartsLocked(keys, starts)
}

func (s *MemoryStorage) XReadManyBlocked(
	keys []string,
	ids []string,
	timeout time.Duration,
) (map[string][]value.StreamEntry, error) {
	if len(keys) != len(ids) {
		return nil, rerrors.ErrInvalidArgs
	}

	deadline := time.Time{}
	if timeout > 0 {
		deadline = time.Now().Add(timeout)
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
	if res != nil {
		s.mu.Unlock()
		return res, nil
	}

	w := &xreadWaiter{ch: make(chan struct{}, 1)}
	for _, key := range keys {
		s.addStreamWaiterLocked(key, w)
	}
	s.mu.Unlock()

	cleanup := func() {
		s.mu.Lock()
		for _, key := range keys {
			s.removeStreamWaiterLocked(key, w)
		}
		s.mu.Unlock()
	}

	if timeout == 0 {
		for {
			<-w.ch

			s.mu.Lock()
			res, err = s.xreadManyFromStartsLocked(keys, starts)
			s.mu.Unlock()

			if err != nil {
				cleanup()
				return nil, err
			}
			if res != nil {
				cleanup()
				return res, nil
			}
		}
	}

	timer := time.NewTimer(time.Until(deadline))
	defer timer.Stop()

	for {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			cleanup()
			return nil, nil
		}

		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
		timer.Reset(remaining)

		select {
		case <-w.ch:
		case <-timer.C:
			cleanup()
			return nil, nil
		}

		s.mu.Lock()
		res, err = s.xreadManyFromStartsLocked(keys, starts)
		s.mu.Unlock()

		if err != nil {
			cleanup()
			return nil, err
		}
		if res != nil {
			cleanup()
			return res, nil
		}
	}
}

package storage

import (
	"errors"
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/ports"
)

type item struct {
	kind valueType

	value string

	list []string

	expiresAt *time.Time
}

var ErrWrongType = errors.New(
	"WRONGTYPE Operation against a key holding the wrong kind of value",
)

type MemoryStorage struct {
	mu   sync.RWMutex
	data map[string]item
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string]item),
	}
}

func (s *MemoryStorage) Set(key string, value string, opts ports.SetOptions) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var expiresAt *time.Time
	if opts.PX != nil {
		t := time.Now().Add(*opts.PX)
		expiresAt = &t
	}

	s.data[key] = item{
		kind:      typeString,
		value:     value,
		expiresAt: expiresAt,
	}

	return nil
}

func (s *MemoryStorage) Get(key string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	it, ok := s.data[key]
	if !ok {
		return "", false, nil
	}

	if it.expiresAt != nil && time.Now().After(*it.expiresAt) {
		delete(s.data, key)
		return "", false, nil
	}

	if it.kind != typeString {
		return "", false, ErrWrongType
	}

	return it.value, true, nil
}

func (s *MemoryStorage) RPush(key string, values ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	it, ok := s.data[key]

	if ok {
		if it.expiresAt != nil && time.Now().After(*it.expiresAt) {
			delete(s.data, key)
			ok = false
		}
	}

	if !ok {
		list := append([]string{}, values...)
		s.data[key] = item{
			kind: typeList,
			list: list,
		}
		return len(list), nil
	}

	if it.kind != typeList {
		return 0, ErrWrongType
	}

	it.list = append(it.list, values...)
	s.data[key] = it

	return len(it.list), nil
}

func (s *MemoryStorage) LPush(key string, values ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	it, ok := s.data[key]

	if ok {
		if it.expiresAt != nil && time.Now().After(*it.expiresAt) {
			delete(s.data, key)
			ok = false
		}
	}

	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}

	if !ok {
		s.data[key] = item{
			kind: typeList,
			list: append([]string{}, values...),
		}
		return len(values), nil
	}

	if it.kind != typeList {
		return 0, ErrWrongType
	}

	it.list = append(values, it.list...)
	s.data[key] = it

	return len(it.list), nil
}

func (s *MemoryStorage) LPop(key string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	it, ok := s.data[key]
	if !ok || it.kind != typeList || len(it.list) == 0 {
		return "", false, nil
	}

	v := it.list[0]
	it.list = it.list[1:]

	if len(it.list) == 0 {
		delete(s.data, key)
	} else {
		s.data[key] = it
	}

	return v, true, nil
}

func (s *MemoryStorage) RPop(key string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	it, ok := s.data[key]
	if !ok || it.kind != typeList || len(it.list) == 0 {
		return "", false, nil
	}

	v := it.list[len(it.list)-1]
	it.list = it.list[:len(it.list)-1]

	if len(it.list) == 0 {
		delete(s.data, key)
	} else {
		s.data[key] = it
	}

	return v, true, nil
}

func (s *MemoryStorage) LRange(key string, start, stop int) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	it, ok := s.data[key]
	if !ok || it.kind != typeList {
		return []string{}, nil
	}

	n := len(it.list)
	if n == 0 {
		return []string{}, nil
	}

	if start < 0 {
		start = n + start
	}
	if stop < 0 {
		stop = n + stop
	}

	if start < 0 {
		start = 0
	}
	if stop >= n {
		stop = n - 1
	}

	if start > stop || start >= n {
		return []string{}, nil
	}

	result := it.list[start : stop+1]

	out := make([]string, len(result))
	copy(out, result)

	return out, nil
}

func (s *MemoryStorage) LLen(key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	it, ok := s.data[key]
	if !ok {
		return 0, nil
	}
	
	if it.expiresAt != nil && time.Now().After(*it.expiresAt) {
		delete(s.data, key)
		return 0, nil
	}

	if it.kind != typeList {
		return 0, ErrWrongType
	}

	return len(it.list), nil
}

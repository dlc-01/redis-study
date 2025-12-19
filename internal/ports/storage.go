package ports

import "time"

type SetOptions struct {
	PX *time.Duration
}

type Storage interface {
	Set(key string, value string, opts SetOptions) error
	Get(key string) (value string, ok bool, err error)

	RPush(key string, values ...string) (int, error)
	LPush(key string, values ...string) (int, error)
	LPop(key string) (string, bool, error)
	RPop(key string) (string, bool, error)
	LRange(key string, start, stop int) ([]string, error)
	LLen(key string) (int, error)
}

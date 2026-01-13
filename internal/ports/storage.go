package ports

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/domain/value"
)

type SetOptions struct {
	PX *time.Duration
}

type Storage interface {
	Set(key string, value string, opts SetOptions) error
	Get(key string) (value string, ok bool, err error)

	RPush(key string, values ...string) (int, error)
	LPush(key string, values ...string) (int, error)
	LPop(key string) (string, bool, error)
	LPopN(key string, count int) ([]string, error)
	RPop(key string) (string, bool, error)
	LRange(key string, start, stop int) ([]string, error)
	LLen(key string) (int, error)
	BLPop(key string) (string, bool, error)
	BLPopWithTimeout(key string, timeout time.Duration) (string, bool, error)
	Type(key string) (string, error)
	XAdd(key string, id string, values []value.StreamKV) (string, error)
	XRange(key string, start, end string) ([]value.StreamEntry, error)
	XReadMany(keys []string, ids []string) (map[string][]value.StreamEntry, error)
}

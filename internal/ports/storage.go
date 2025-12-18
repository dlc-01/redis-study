package ports

import "time"

type SetOptions struct {
	PX *time.Duration
}

type Storage interface {
	Set(key string, value string, opts SetOptions) error
	Get(key string) (value string, ok bool, err error)
}

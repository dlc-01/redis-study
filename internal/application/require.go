package application

import "github.com/codecrafters-io/redis-starter-go/internal/domain/rerrors"

func require[T any](cmd Command) (T, error) {
	v, ok := cmd.(T)
	if !ok {
		var zero T
		return zero, rerrors.ErrInternal
	}
	return v, nil
}

package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/internal/application/commands/list"
	"github.com/codecrafters-io/redis-starter-go/internal/application/errors"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/response"
	"github.com/codecrafters-io/redis-starter-go/internal/ports"
)

func HandleLLen(storage ports.Storage, c list.LLenCommand) (response.Response, error) {
	n, err := storage.LLen(c.Key)
	if err != nil {
		return errors.ToResponse(err), nil
	}
	return response.Integer{Value: n}, nil
}

package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/internal/application/commands/list"
	"github.com/codecrafters-io/redis-starter-go/internal/application/errors"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/rerrors"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/response"
	"github.com/codecrafters-io/redis-starter-go/internal/ports"
)

func HandleRPush(storage ports.Storage, c list.RPushCommand) (response.Response, error) {
	if len(c.Values) == 0 {
		return response.ErrorString{Message: rerrors.WrongNumberOfArgs("rpush").Error()}, nil

	}

	n, err := storage.RPush(c.Key, c.Values...)
	if err != nil {
		return errors.ToResponse(err), nil
	}

	return response.Integer{Value: n}, nil
}

func HandleLPush(storage ports.Storage, c list.LPushCommand) (response.Response, error) {
	if len(c.Values) == 0 {
		return response.ErrorString{Message: rerrors.WrongNumberOfArgs("lpush").Error()}, nil
	}

	n, err := storage.LPush(c.Key, c.Values...)
	if err != nil {
		return errors.ToResponse(err), nil
	}

	return response.Integer{Value: n}, nil
}

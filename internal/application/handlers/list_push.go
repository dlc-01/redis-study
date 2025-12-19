package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/internal/application/commands/list"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/response"
	"github.com/codecrafters-io/redis-starter-go/internal/ports"
)

func HandleRPush(storage ports.Storage, c list.RPushCommand) (response.Response, error) {
	if len(c.Values) == 0 {
		return response.ErrorString{
			Message: "ERR wrong number of arguments for 'rpush' command",
		}, nil
	}

	n, err := storage.RPush(c.Key, c.Values...)
	if err != nil {
		return response.ErrorString{Message: err.Error()}, nil
	}

	return response.Integer{Value: n}, nil
}

func HandleLPush(storage ports.Storage, c list.LPushCommand) (response.Response, error) {
	if len(c.Values) == 0 {
		return response.ErrorString{
			Message: "ERR wrong number of arguments for 'lpush' command",
		}, nil
	}

	n, err := storage.LPush(c.Key, c.Values...)
	if err != nil {
		return response.ErrorString{Message: err.Error()}, nil
	}

	return response.Integer{Value: n}, nil
}

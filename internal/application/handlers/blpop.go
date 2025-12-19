package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/internal/application/commands/list"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/response"
	"github.com/codecrafters-io/redis-starter-go/internal/ports"
)

func HandleBLPop(storage ports.Storage, c list.BLPopCommand) (response.Response, error) {
	var (
		val string
		ok  bool
		err error
	)

	if c.Timeout == 0 {
		val, ok, err = storage.BLPop(c.Key)
	} else {
		val, ok, err = storage.BLPopWithTimeout(c.Key, c.Timeout)
	}

	if err != nil {
		return response.ErrorString{Message: err.Error()}, nil
	}

	if !ok {
		return response.NullArray{}, nil
	}

	key := c.Key
	value := val

	return response.Array{
		Values: []response.Response{
			response.BulkString{Value: &key},
			response.BulkString{Value: &value},
		},
	}, nil
}

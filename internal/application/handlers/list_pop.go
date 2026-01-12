package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/internal/application/commands/list"
	"github.com/codecrafters-io/redis-starter-go/internal/application/errors"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/response"
	"github.com/codecrafters-io/redis-starter-go/internal/ports"
)

func HandleLPop(storage ports.Storage, c list.LPopCommand) (response.Response, error) {
	if c.Count == nil {
		v, ok, err := storage.LPop(c.Key)
		if err != nil {
			return errors.ToResponse(err), nil
		}
		if !ok {
			return response.BulkString{Value: nil}, nil
		}
		return response.BulkString{Value: &v}, nil
	}

	values, err := storage.LPopN(c.Key, *c.Count)
	if err != nil {
		return errors.ToResponse(err), nil
	}

	resp := make([]response.Response, 0, len(values))
	for _, v := range values {
		val := v
		resp = append(resp, response.BulkString{Value: &val})
	}

	return response.Array{Values: resp}, nil
}

package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/internal/application/commands/list"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/response"
	"github.com/codecrafters-io/redis-starter-go/internal/ports"
)

func HandleLRange(storage ports.Storage, c list.LRangeCommand) (response.Response, error) {
	values, err := storage.LRange(c.Key, c.Start, c.Stop)
	if err != nil {
		return response.ErrorString{Message: "ERR internal error"}, nil
	}

	resp := make([]response.Response, 0, len(values))
	for _, v := range values {
		val := v
		resp = append(resp, response.BulkString{Value: &val})
	}

	return response.Array{Values: resp}, nil
}

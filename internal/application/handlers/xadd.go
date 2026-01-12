package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/internal/application/commands/stream"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/response"
	"github.com/codecrafters-io/redis-starter-go/internal/ports"
)

func HandleXAdd(storage ports.Storage, c stream.XAddCommand) (response.Response, error) {
	id, err := storage.XAdd(c.Key, c.ID, c.Values)
	if err != nil {
		return nil, err
	}

	return response.BulkString{Value: &id}, nil
}

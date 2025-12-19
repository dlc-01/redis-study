package handlers

import (
	commandString "github.com/codecrafters-io/redis-starter-go/internal/application/commands/string"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/response"
	"github.com/codecrafters-io/redis-starter-go/internal/ports"
)

func HandleSet(storage ports.Storage, c commandString.SetCommand) (response.Response, error) {
	if err := storage.Set(c.Key, c.Value, ports.SetOptions{PX: c.PX}); err != nil {
		return response.ErrorString{Message: "ERR internal error"}, nil
	}
	return response.SimpleString{Value: "OK"}, nil
}

func HandleGet(storage ports.Storage, c commandString.GetCommand) (response.Response, error) {
	val, ok, err := storage.Get(c.Key)
	if err != nil {
		return response.ErrorString{Message: "ERR internal error"}, nil
	}
	if !ok {
		return response.BulkString{Value: nil}, nil
	}
	return response.BulkString{Value: &val}, nil
}

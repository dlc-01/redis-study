package handlers

import (
	commandString "github.com/codecrafters-io/redis-starter-go/internal/application/commands/string"
	"github.com/codecrafters-io/redis-starter-go/internal/application/errors"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/response"
	"github.com/codecrafters-io/redis-starter-go/internal/ports"
)

func HandleType(storage ports.Storage, c commandString.TypeCommand) (response.Response, error) {
	t, err := storage.Type(c.Key)
	if err != nil {
		return errors.ToResponse(err), nil
	}

	return response.SimpleString{Value: t}, nil
}

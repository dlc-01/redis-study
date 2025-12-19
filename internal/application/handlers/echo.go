package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/internal/application/commands"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/response"
)

func HandleEcho(c commands.EchoCommand) (response.Response, error) {
	return response.BulkString{Value: &c.Value}, nil
}

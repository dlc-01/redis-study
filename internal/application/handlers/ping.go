package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/internal/application/commands"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/response"
)

func HandlePing(_ commands.PingCommand) (response.Response, error) {
	return response.SimpleString{Value: "PONG"}, nil
}

package application

import (
	"github.com/codecrafters-io/redis-starter-go/internal/application/commands"
	"github.com/codecrafters-io/redis-starter-go/internal/application/commands/list"
	commandString "github.com/codecrafters-io/redis-starter-go/internal/application/commands/string"
	"github.com/codecrafters-io/redis-starter-go/internal/application/handlers"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/response"
	"github.com/codecrafters-io/redis-starter-go/internal/ports"
)

type Processor struct {
	storage ports.Storage
}

func NewProcessor(storage ports.Storage) *Processor {
	return &Processor{storage: storage}
}

func (p *Processor) Handle(cmd Command) (response.Response, error) {
	switch c := cmd.(type) {

	case commands.PingCommand:
		return handlers.HandlePing(c)

	case commands.EchoCommand:
		return handlers.HandleEcho(c)

	case commandString.SetCommand:
		return handlers.HandleSet(p.storage, c)

	case commandString.GetCommand:
		return handlers.HandleGet(p.storage, c)

	case list.RPushCommand:
		return handlers.HandleRPush(p.storage, c)

	case list.LPushCommand:
		return handlers.HandleLPush(p.storage, c)

	case list.LRangeCommand:
		return handlers.HandleLRange(p.storage, c)

	case list.LLenCommand:
		return handlers.HandleLLen(p.storage, c)

	case list.LPopCommand:
		return handlers.HandleLPop(p.storage, c)

	case list.BLPopCommand:
		return handlers.HandleBLPop(p.storage, c)

	case commandString.TypeCommand:
		return handlers.HandleType(p.storage, c)
	}

	return response.ErrorString{Message: "ERR unknown command"}, nil
}

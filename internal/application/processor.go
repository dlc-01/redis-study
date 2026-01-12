package application

import (
	"github.com/codecrafters-io/redis-starter-go/internal/application/commands"
	"github.com/codecrafters-io/redis-starter-go/internal/application/commands/list"
	"github.com/codecrafters-io/redis-starter-go/internal/application/commands/stream"
	commandString "github.com/codecrafters-io/redis-starter-go/internal/application/commands/string"
	apperrors "github.com/codecrafters-io/redis-starter-go/internal/application/errors"
	"github.com/codecrafters-io/redis-starter-go/internal/application/handlers"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/response"
	"github.com/codecrafters-io/redis-starter-go/internal/ports"
)

type Processor struct {
	router *Router
}

func NewProcessor(storage ports.Storage) *Processor {
	r := NewRouter()

	r.Register("PING", func(cmd Command) (response.Response, error) {
		c, err := require[commands.PingCommand](cmd)
		if err != nil {
			return nil, err
		}
		return handlers.HandlePing(c)
	})

	r.Register("ECHO", func(cmd Command) (response.Response, error) {
		c, err := require[commands.EchoCommand](cmd)
		if err != nil {
			return nil, err
		}
		return handlers.HandleEcho(c)
	})

	r.Register("SET", func(cmd Command) (response.Response, error) {
		c, err := require[commandString.SetCommand](cmd)
		if err != nil {
			return nil, err
		}
		return handlers.HandleSet(storage, c)
	})

	r.Register("GET", func(cmd Command) (response.Response, error) {
		c, err := require[commandString.GetCommand](cmd)
		if err != nil {
			return nil, err
		}
		return handlers.HandleGet(storage, c)
	})

	r.Register("TYPE", func(cmd Command) (response.Response, error) {
		c, err := require[commandString.TypeCommand](cmd)
		if err != nil {
			return nil, err
		}
		return handlers.HandleType(storage, c)
	})

	r.Register("RPUSH", func(cmd Command) (response.Response, error) {
		c, err := require[list.RPushCommand](cmd)
		if err != nil {
			return nil, err
		}
		return handlers.HandleRPush(storage, c)
	})

	r.Register("LPUSH", func(cmd Command) (response.Response, error) {
		c, err := require[list.LPushCommand](cmd)
		if err != nil {
			return nil, err
		}
		return handlers.HandleLPush(storage, c)
	})

	r.Register("LRANGE", func(cmd Command) (response.Response, error) {
		c, err := require[list.LRangeCommand](cmd)
		if err != nil {
			return nil, err
		}
		return handlers.HandleLRange(storage, c)
	})

	r.Register("LLEN", func(cmd Command) (response.Response, error) {
		c, err := require[list.LLenCommand](cmd)
		if err != nil {
			return nil, err
		}
		return handlers.HandleLLen(storage, c)
	})

	r.Register("LPOP", func(cmd Command) (response.Response, error) {
		c, err := require[list.LPopCommand](cmd)
		if err != nil {
			return nil, err
		}
		return handlers.HandleLPop(storage, c)
	})

	r.Register("BLPOP", func(cmd Command) (response.Response, error) {
		c, err := require[list.BLPopCommand](cmd)
		if err != nil {
			return nil, err
		}
		return handlers.HandleBLPop(storage, c)
	})

	r.Register("XADD", func(cmd Command) (response.Response, error) {
		c, err := require[stream.XAddCommand](cmd)
		if err != nil {
			return nil, err
		}
		return handlers.HandleXAdd(storage, c)
	})

	r.Register("XRANGE", func(cmd Command) (response.Response, error) {
		c, err := require[stream.XRangeCommand](cmd)
		if err != nil {
			return nil, err
		}
		return handlers.HandleXRange(storage, c)
	})

	return &Processor{router: r}
}

func (p *Processor) Handle(cmd Command) (response.Response, error) {
	resp, err := p.router.Handle(cmd)
	if err != nil {
		return apperrors.ToResponse(err), nil
	}
	return resp, nil
}

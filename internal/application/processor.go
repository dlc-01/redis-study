package application

import (
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

	case PingCommand:
		return response.SimpleString{Value: "PONG"}, nil

	case EchoCommand:
		return response.BulkString{Value: &c.Value}, nil

	case SetCommand:
		if err := p.storage.Set(c.Key, c.Value); err != nil {
			return response.ErrorString{Message: "ERR internal error"}, nil
		}
		return response.SimpleString{Value: "OK"}, nil

	case GetCommand:
		val, ok, err := p.storage.Get(c.Key)
		if err != nil {
			return response.ErrorString{Message: "ERR internal error"}, nil
		}
		if !ok {
			return response.BulkString{Value: nil}, nil
		}
		return response.BulkString{Value: &val}, nil
	}

	return response.ErrorString{Message: "ERR unknown command"}, nil
}

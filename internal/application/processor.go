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
		if err := p.storage.Set(c.Key, c.Value, ports.SetOptions{PX: c.PX}); err != nil {
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

	case RPushCommand:
		if len(c.Values) == 0 {
			return response.ErrorString{
				Message: "ERR wrong number of arguments for 'rpush' command",
			}, nil
		}

		n, err := p.storage.RPush(c.Key, c.Values...)
		if err != nil {
			return response.ErrorString{Message: err.Error()}, nil
		}

		return response.Integer{Value: n}, nil
	case LPushCommand:
		if len(c.Values) == 0 {
			return response.ErrorString{
				Message: "ERR wrong number of arguments for 'lpush' command",
			}, nil
		}

		n, err := p.storage.LPush(c.Key, c.Values...)
		if err != nil {
			return response.ErrorString{Message: err.Error()}, nil
		}

		return response.Integer{Value: n}, nil

	case LRangeCommand:
		values, err := p.storage.LRange(c.Key, c.Start, c.Stop)
		if err != nil {
			return response.ErrorString{Message: "ERR internal error"}, nil
		}

		resp := make([]response.Response, 0, len(values))
		for _, v := range values {
			val := v
			resp = append(resp, response.BulkString{Value: &val})
		}

		return response.Array{Values: resp}, nil
	case LLenCommand:
		n, err := p.storage.LLen(c.Key)
		if err != nil {
			return response.ErrorString{Message: err.Error()}, nil
		}
		return response.Integer{Value: n}, nil

	case LPopCommand:
		if c.Count == nil {
			v, ok, err := p.storage.LPop(c.Key)
			if err != nil {
				return response.ErrorString{Message: err.Error()}, nil
			}
			if !ok {
				return response.BulkString{Value: nil}, nil
			}
			return response.BulkString{Value: &v}, nil
		}

		values, err := p.storage.LPopN(c.Key, *c.Count)
		if err != nil {
			return response.ErrorString{Message: err.Error()}, nil
		}

		resp := make([]response.Response, 0, len(values))
		for _, v := range values {
			val := v
			resp = append(resp, response.BulkString{Value: &val})
		}

		return response.Array{Values: resp}, nil

	}

	return response.ErrorString{Message: "ERR unknown command"}, nil
}

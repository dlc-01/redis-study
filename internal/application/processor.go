package application

import (
	"errors"

	"github.com/codecrafters-io/redis-starter-go/internal/domain/response"
)

var (
	ErrUnknownCommand = errors.New("unknown command")
	ErrWrongArity     = errors.New("wrong number of arguments")
)

type Processor struct{}

func NewProcessor() *Processor {
	return &Processor{}
}

func (p *Processor) Process(cmd Command) (response.Response, error) {
	switch cmd.Name {

	case "PING":
		return response.SimpleString{Value: "PONG"}, nil

	case "ECHO":
		if len(cmd.Args) != 1 {
			return nil, ErrWrongArity
		}
		return response.BulkString{Value: cmd.Args[0]}, nil

	default:
		return nil, ErrUnknownCommand
	}
}

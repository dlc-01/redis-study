package parser

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/application"
	"github.com/codecrafters-io/redis-starter-go/internal/application/commands/stream"
)

func (p *RespParser) parseStream(args []string) (application.Command, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("ERR wrong number of arguments for 'xadd' command")
	}

	key := args[1]
	id := args[2]

	rest := args[3:]
	if len(rest)%2 != 0 {
		return nil, fmt.Errorf("ERR wrong number of arguments for 'xadd' command")
	}

	fields := make(map[string]string, len(rest)/2)
	for i := 0; i < len(rest); i += 2 {
		k := rest[i]
		v := rest[i+1]
		fields[k] = v
	}

	return stream.XAddCommand{
		Key:    key,
		ID:     id,
		Fields: fields,
	}, nil
}

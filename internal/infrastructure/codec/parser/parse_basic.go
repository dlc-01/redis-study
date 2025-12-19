package parser

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/application"
	"github.com/codecrafters-io/redis-starter-go/internal/application/commands"
)

func (p *RespParser) parseBasic(args []string) (application.Command, error) {
	switch strings.ToUpper(args[0]) {

	case "PING":
		return commands.PingCommand{}, nil

	case "ECHO":
		if len(args) != 2 {
			return nil, fmt.Errorf("ECHO expects 1 argument")
		}
		return commands.EchoCommand{Value: args[1]}, nil
	}

	return nil, fmt.Errorf("unknown basic command")
}

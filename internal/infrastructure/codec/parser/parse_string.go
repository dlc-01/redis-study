package parser

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/application"
	commandString "github.com/codecrafters-io/redis-starter-go/internal/application/commands/string"
)

func (p *RespParser) parseString(args []string) (application.Command, error) {
	switch strings.ToUpper(args[0]) {

	case "GET":
		if len(args) != 2 {
			return nil, fmt.Errorf("GET expects 1 argument")
		}
		return commandString.GetCommand{Key: args[1]}, nil

	case "SET":
		return p.parseSet(args)
	}

	return nil, fmt.Errorf("unknown string command")
}

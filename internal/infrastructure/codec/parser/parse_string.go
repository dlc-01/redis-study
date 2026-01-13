package parser

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/application"
	commandString "github.com/codecrafters-io/redis-starter-go/internal/application/commands/string"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/rerrors"
)

func (p *RespParser) parseString(args []string) (application.Command, error) {
	switch strings.ToUpper(args[0]) {

	case "GET":
		if len(args) != 2 {
			return nil, rerrors.WrongNumberOfArgs("get")
		}
		return commandString.GetCommand{Key: args[1]}, nil

	case "SET":
		return p.parseSet(args)

	case "TYPE":
		if len(args) != 2 {
			return nil, rerrors.WrongNumberOfArgs("type")
		}
		return commandString.TypeCommand{
			Key: args[1],
		}, nil

	}

	return nil, fmt.Errorf("unknown string command")
}

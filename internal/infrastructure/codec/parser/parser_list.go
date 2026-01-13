package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/application"
	"github.com/codecrafters-io/redis-starter-go/internal/application/commands/list"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/rerrors"
)

func (p *RespParser) parseList(args []string) (application.Command, error) {
	switch strings.ToUpper(args[0]) {

	case "RPUSH":
		if len(args) < 3 {
			return nil, rerrors.WrongNumberOfArgs("rpush")
		}
		return list.RPushCommand{Key: args[1], Values: args[2:]}, nil

	case "LPUSH":
		if len(args) < 3 {
			return nil, rerrors.WrongNumberOfArgs("lpush")
		}
		return list.LPushCommand{Key: args[1], Values: args[2:]}, nil

	case "LLEN":
		if len(args) != 2 {
			return nil, rerrors.WrongNumberOfArgs("llen")
		}
		return list.LLenCommand{Key: args[1]}, nil

	case "LRANGE":
		if len(args) != 4 {
			return nil, rerrors.WrongNumberOfArgs("lrange")
		}
		start, err := strconv.Atoi(args[2])
		if err != nil {
			return nil, fmt.Errorf("invalid start index")
		}
		stop, err := strconv.Atoi(args[3])
		if err != nil {
			return nil, fmt.Errorf("invalid stop index")
		}
		return list.LRangeCommand{Key: args[1], Start: start, Stop: stop}, nil

	case "LPOP":
		if len(args) == 2 {
			return list.LPopCommand{Key: args[1]}, nil
		}
		if len(args) == 3 {
			n, err := strconv.Atoi(args[2])
			if err != nil || n < 0 {
				return nil, fmt.Errorf("invalid count")
			}
			return list.LPopCommand{Key: args[1], Count: &n}, nil
		}
		return nil, rerrors.WrongNumberOfArgs("lpop")

	case "BLPOP":
		if len(args) != 3 {
			return nil, rerrors.WrongNumberOfArgs("blpop")
		}

		timeoutSec, err := strconv.ParseFloat(args[2], 64)
		if err != nil || timeoutSec < 0 {
			return nil, fmt.Errorf("invalid timeout")
		}

		timeout := time.Duration(timeoutSec * float64(time.Second))

		return list.BLPopCommand{
			Key:     args[1],
			Timeout: timeout,
		}, nil
	}

	return nil, fmt.Errorf("unknown list command")
}

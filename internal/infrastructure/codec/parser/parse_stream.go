package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/application"
	"github.com/codecrafters-io/redis-starter-go/internal/application/commands/stream"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/rerrors"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/value"
)

func (p *RespParser) parseStream(args []string) (application.Command, error) {
	switch strings.ToUpper(args[0]) {

	case "XADD":
		if len(args) < 5 {
			return nil, rerrors.WrongNumberOfArgs("xadd")
		}

		key := args[1]
		id := args[2]

		rest := args[3:]
		if len(rest)%2 != 0 {
			return nil, rerrors.WrongNumberOfArgs("xadd")
		}

		kvs := make([]value.StreamKV, 0, len(rest)/2)
		for i := 0; i < len(rest); i += 2 {
			kvs = append(kvs, value.StreamKV{Key: rest[i], Value: rest[i+1]})
		}

		return stream.XAddCommand{Key: key, ID: id, Values: kvs}, nil

	case "XRANGE":
		if len(args) != 4 {
			return nil, rerrors.WrongNumberOfArgs("xrange")
		}
		return stream.XRangeCommand{
			Key:   args[1],
			Start: args[2],
			End:   args[3],
		}, nil

	case "XREAD":
		if len(args) < 4 {
			return nil, rerrors.WrongNumberOfArgs("xread")
		}

		i := 1
		var block *time.Duration

		if strings.ToUpper(args[i]) == "BLOCK" {
			if len(args) < 6 {
				return nil, rerrors.WrongNumberOfArgs("xread")
			}
			ms, err := strconv.ParseInt(args[i+1], 10, 64)
			if err != nil || ms < 0 {
				return nil, fmt.Errorf("ERR invalid BLOCK value")
			}
			d := time.Duration(ms) * time.Millisecond
			block = &d
			i += 2
		}

		if strings.ToUpper(args[i]) != "STREAMS" {
			return nil, fmt.Errorf("ERR syntax error")
		}
		i++

		rest := args[i:]
		if len(rest)%2 != 0 || len(rest) == 0 {
			return nil, fmt.Errorf("ERR syntax error")
		}

		n := len(rest) / 2
		keys := append([]string{}, rest[:n]...)
		ids := append([]string{}, rest[n:]...)

		return stream.XReadCommand{Block: block, Keys: keys, IDs: ids}, nil
	}
	return nil, fmt.Errorf("unknown command")
}

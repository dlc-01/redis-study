package codec

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/application"
)

type RespParser struct {
	reader *bufio.Reader
}

func NewRespParser(r *bufio.Reader) *RespParser {
	return &RespParser{reader: r}
}

func (p *RespParser) ReadCommand() (application.Command, error) {
	args, err := p.readArray()
	if err != nil {
		return nil, err
	}

	if len(args) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	switch strings.ToUpper(args[0]) {

	case "PING":
		return application.PingCommand{}, nil

	case "ECHO":
		if len(args) != 2 {
			return nil, fmt.Errorf("ECHO expects 1 argument")
		}
		return application.EchoCommand{Value: args[1]}, nil

	case "SET":
		return p.parseSet(args)

	case "GET":
		if len(args) != 2 {
			return nil, fmt.Errorf("GET expects 1 argument")
		}
		return application.GetCommand{
			Key: args[1],
		}, nil

	case "RPUSH":
		if len(args) < 3 {
			return nil, fmt.Errorf("RPUSH expects at least 2 arguments")
		}
		return application.RPushCommand{
			Key:    args[1],
			Values: args[2:],
		}, nil
	case "LPUSH":
		if len(args) < 3 {
			return nil, fmt.Errorf("LPUSH expects at least 2 arguments")
		}
		return application.LPushCommand{
			Key:    args[1],
			Values: args[2:],
		}, nil
	case "LRANGE":
		if len(args) != 4 {
			return nil, fmt.Errorf("LRANGE expects 3 arguments")
		}

		start, err := strconv.Atoi(args[2])
		if err != nil {
			return nil, fmt.Errorf("invalid start index")
		}

		stop, err := strconv.Atoi(args[3])
		if err != nil {
			return nil, fmt.Errorf("invalid stop index")
		}

		return application.LRangeCommand{
			Key:   args[1],
			Start: start,
			Stop:  stop,
		}, nil

	}

	return nil, fmt.Errorf("unknown command")
}

func (p *RespParser) readArray() ([]string, error) {
	b, err := p.reader.ReadByte()
	if err != nil {
		return nil, err
	}
	if b != '*' {
		return nil, fmt.Errorf("expected array, got %q", b)
	}

	n, err := p.readInt()
	if err != nil {
		return nil, err
	}

	if n < 0 {
		return nil, fmt.Errorf("invalid array size")
	}

	result := make([]string, 0, n)

	for i := 0; i < n; i++ {
		s, err := p.readBulkString()
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}

	return result, nil
}

func (p *RespParser) readInt() (int, error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	line = line[:len(line)-2]

	n, err := strconv.Atoi(line)
	if err != nil {
		return 0, fmt.Errorf("invalid integer: %q", line)
	}

	return n, nil
}

func (p *RespParser) readBulkString() (string, error) {
	b, err := p.reader.ReadByte()
	if err != nil {
		return "", err
	}
	if b != '$' {
		return "", fmt.Errorf("expected bulk string, got %q", b)
	}

	n, err := p.readInt()
	if err != nil {
		return "", err
	}

	if n < 0 {
		return "", nil
	}

	buf := make([]byte, n)
	_, err = p.reader.Read(buf)
	if err != nil {
		return "", err
	}

	cr, err := p.reader.ReadByte()
	if err != nil || cr != '\r' {
		return "", fmt.Errorf("expected CR after bulk string")
	}

	lf, err := p.reader.ReadByte()
	if err != nil || lf != '\n' {
		return "", fmt.Errorf("expected LF after bulk string")
	}

	return string(buf), nil
}

func (p *RespParser) parseSet(args []string) (application.Command, error) {

	if len(args) < 3 {
		return nil, fmt.Errorf("SET expects at least 2 arguments")
	}

	key := args[1]
	value := args[2]

	var px *time.Duration

	i := 3
	for i < len(args) {
		switch strings.ToUpper(args[i]) {

		case "PX":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("PX requires a value")
			}

			ms, err := strconv.Atoi(args[i+1])
			if err != nil || ms < 0 {
				return nil, fmt.Errorf("invalid PX value")
			}

			d := time.Duration(ms) * time.Millisecond
			px = &d
			i += 2

		default:
			return nil, fmt.Errorf("unknown SET option: %s", args[i])
		}
	}

	return application.SetCommand{
		Key:   key,
		Value: value,
		PX:    px,
	}, nil
}

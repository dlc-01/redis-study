package codec

import (
	"bufio"
	"errors"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/application"
)

type RespParser struct {
	reader *bufio.Reader
}

func NewRespParser(r *bufio.Reader) *RespParser {
	return &RespParser{reader: r}
}

func (p *RespParser) ReadCommand() (application.Command, error) {
	size, err := p.readArraySize()
	if err != nil {
		return application.Command{}, err
	}

	parts := make([]string, 0, size)
	for i := 0; i < size; i++ {
		s, err := p.readBulkString()
		if err != nil {
			return application.Command{}, err
		}
		parts = append(parts, s)
	}

	if len(parts) == 0 {
		return application.Command{}, errors.New("empty command")
	}

	return application.Command{
		Name: strings.ToUpper(parts[0]),
		Args: parts[1:],
	}, nil
}

func (p *RespParser) readArraySize() (int, error) {
	b, err := p.reader.ReadByte()
	if err != nil || b != '*' {
		return 0, errors.New("expected array")
	}

	line, err := p.reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(strings.TrimSpace(line))
}

func (p *RespParser) readBulkString() (string, error) {
	b, err := p.reader.ReadByte()
	if err != nil || b != '$' {
		return "", errors.New("expected bulk string")
	}

	line, err := p.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	n, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		return "", err
	}

	buf := make([]byte, n)
	if _, err := p.reader.Read(buf); err != nil {
		return "", err
	}
	
	p.reader.ReadByte()
	p.reader.ReadByte()

	return string(buf), nil
}

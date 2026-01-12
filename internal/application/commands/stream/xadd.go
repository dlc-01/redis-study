package stream

import "github.com/codecrafters-io/redis-starter-go/internal/domain/value"

type XAddCommand struct {
	Key    string
	ID     string
	Values []value.StreamKV
}

func (XAddCommand) Name() string { return "XADD" }

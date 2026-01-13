package stream

import "time"

type XReadCommand struct {
	Block *time.Duration
	Keys  []string
	IDs   []string
}

func (XReadCommand) Name() string { return "XREAD" }

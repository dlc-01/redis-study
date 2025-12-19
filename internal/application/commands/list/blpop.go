package list

import "time"

type BLPopCommand struct {
	Key     string
	Timeout time.Duration
}

func (BLPopCommand) Name() string {
	return "BLPOP"
}

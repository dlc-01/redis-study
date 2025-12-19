package string

import "time"

type SetCommand struct {
	Key   string
	Value string
	PX    *time.Duration
}

func (SetCommand) Name() string {
	return "SET"
}

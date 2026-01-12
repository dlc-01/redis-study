package rerrors

import (
	"errors"
	"fmt"
)

var (
	ErrWrongType   = errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	ErrInvalidArgs = errors.New("ERR wrong number of arguments")
	ErrUnknown     = errors.New("ERR unknown command")
	ErrInternal    = errors.New("ERR internal error")
)

func WrongNumberOfArgs(cmd string) error {
	return fmt.Errorf("ERR wrong number of arguments for '%s' command", cmd)
}

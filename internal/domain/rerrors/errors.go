package rerrors

import (
	"errors"
	"fmt"
)

var (
	ErrWrongType       = errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	ErrInvalidArgs     = errors.New("ERR wrong number of arguments")
	ErrUnknown         = errors.New("ERR unknown command")
	ErrInternal        = errors.New("ERR internal error")
	ErrXAddIDTooSmall  = errors.New("ERR The ID specified in XADD is equal or smaller than the target stream top item")
	ErrXAddZeroID      = errors.New("ERR The ID specified in XADD must be greater than 0-0")
	ErrInvalidStreamID = errors.New("ERR Invalid stream ID specified as stream command argument")
)

func WrongNumberOfArgs(cmd string) error {
	return fmt.Errorf("ERR wrong number of arguments for '%s' command", cmd)
}

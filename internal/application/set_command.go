package application

type SetCommand struct {
	Key   string
	Value string
}

func (SetCommand) Name() string {
	return "SET"
}

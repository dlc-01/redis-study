package string

type TypeCommand struct {
	Key string
}

func (TypeCommand) Name() string {
	return "TYPE"
}

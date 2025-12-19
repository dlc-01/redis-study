package string

type GetCommand struct {
	Key string
}

func (GetCommand) Name() string {
	return "GET"
}

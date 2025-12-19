package application

type RPushCommand struct {
	Key    string
	Values []string
}

func (RPushCommand) Name() string {
	return "RPUSH"
}

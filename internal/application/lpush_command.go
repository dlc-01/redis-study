package application

type LPushCommand struct {
	Key    string
	Values []string
}

func (LPushCommand) Name() string {
	return "LPUSH"
}

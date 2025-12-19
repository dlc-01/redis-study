package list

type LPushCommand struct {
	Key    string
	Values []string
}

func (LPushCommand) Name() string {
	return "LPUSH"
}

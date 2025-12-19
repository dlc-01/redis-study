package list

type LLenCommand struct {
	Key string
}

func (LLenCommand) Name() string {
	return "LLEN"
}

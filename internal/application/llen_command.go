package application

type LLenCommand struct {
	Key string
}

func (LLenCommand) Name() string {
	return "LLEN"
}

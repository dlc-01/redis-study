package application

type LPopCommand struct {
	Key   string
	Count *int
}

func (LPopCommand) Name() string {
	return "LPOP"
}

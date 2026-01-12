package stream

type XAddCommand struct {
	Key    string
	ID     string
	Fields map[string]string
}

func (XAddCommand) Name() string {
	return "XADD"
}

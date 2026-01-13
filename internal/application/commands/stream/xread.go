package stream

type XReadCommand struct {
	Keys []string
	IDs  []string
}

func (XReadCommand) Name() string { return "XREAD" }

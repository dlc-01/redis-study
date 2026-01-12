package stream

type XRangeCommand struct {
	Key   string
	Start string
	End   string
}

func (XRangeCommand) Name() string { return "XRANGE" }

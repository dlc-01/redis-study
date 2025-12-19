package list

type LRangeCommand struct {
	Key   string
	Start int
	Stop  int
}

func (LRangeCommand) Name() string {
	return "LRANGE"
}

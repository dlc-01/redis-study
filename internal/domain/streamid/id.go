package streamid

import "fmt"

type ID struct {
	Time int64
	Seq  int64
}

func Compare(a, b ID) int {
	if a.Time < b.Time {
		return -1
	}
	if a.Time > b.Time {
		return 1
	}
	if a.Seq < b.Seq {
		return -1
	}
	if a.Seq > b.Seq {
		return 1
	}
	return 0
}

func (id ID) String() string { return fmt.Sprintf("%d-%d", id.Time, id.Seq) }

package value

type Stream struct {
	V []StreamEntry
}

func (Stream) Type() Type { return TypeStream }

type StreamKV struct {
	Key   string
	Value string
}

type StreamEntry struct {
	ID     string
	Values []StreamKV
}

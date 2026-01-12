package value

type StreamEntry struct {
	ID     string
	Fields map[string]string
}

type Stream struct {
	V []StreamEntry
}

func (Stream) Type() Type {
	return TypeStream
}

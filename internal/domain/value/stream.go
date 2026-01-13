package value

import "github.com/codecrafters-io/redis-starter-go/internal/domain/streamid"

type Stream struct {
	V      []StreamEntry
	LastID streamid.ID
}

func (Stream) Type() Type { return TypeStream }

type StreamKV struct {
	Key   string
	Value string
}

type StreamEntry struct {
	IDRaw  string
	ID     streamid.ID
	Values []StreamKV
}

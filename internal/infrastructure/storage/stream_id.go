package storage

import (
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/domain/rerrors"
)

type streamID struct {
	Time int64
	Seq  int64
}

func parseStreamID(id string) (streamID, error) {
	parts := strings.Split(id, "-")
	if len(parts) != 2 {
		return streamID{}, rerrors.ErrInvalidStreamID
	}

	t, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return streamID{}, err
	}

	s, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return streamID{}, err
	}

	return streamID{Time: t, Seq: s}, nil
}

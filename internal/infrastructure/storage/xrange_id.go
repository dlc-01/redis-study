package storage

import (
	"math"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/domain/rerrors"
)

type rangeID = streamID

func parseRangeID(s string, isStart bool) (streamID, error) {
	if s == "-" {
		return streamID{Time: 0, Seq: 0}, nil
	}
	if s == "+" {
		return streamID{Time: math.MaxInt64, Seq: math.MaxInt64}, nil
	}

	if strings.Contains(s, "-") {
		spec, err := parseIDSpec(s)
		if err != nil {
			return streamID{}, err
		}
		if spec.kind != specExplicit {
			return streamID{}, rerrors.ErrInvalidArgs
		}
		return streamID{Time: spec.t, Seq: spec.seq}, nil
	}

	t, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return streamID{}, rerrors.ErrInvalidArgs
	}

	if isStart {
		return streamID{Time: t, Seq: 0}, nil
	}
	return streamID{Time: t, Seq: math.MaxInt64}, nil
}

func cmpID(a, b streamID) int {
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

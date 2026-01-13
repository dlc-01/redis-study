package streamidcodec

import (
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/domain/rerrors"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/streamid"
)

type SpecKind int

const (
	SpecExplicit SpecKind = iota
	SpecAutoSeq
	SpecAutoID
)

type Spec struct {
	Kind SpecKind
	Time int64
	Seq  int64
}

func ParseSpec(s string) (Spec, error) {
	if s == "*" {
		return Spec{Kind: SpecAutoID}, nil
	}

	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return Spec{}, rerrors.ErrInvalidStreamID
	}

	t, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return Spec{}, rerrors.ErrInvalidStreamID
	}

	if parts[1] == "*" {
		return Spec{Kind: SpecAutoSeq, Time: t}, nil
	}

	seq, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return Spec{}, rerrors.ErrInvalidStreamID
	}

	return Spec{Kind: SpecExplicit, Time: t, Seq: seq}, nil
}

func ParseRangeID(s string, isStart bool) (streamid.ID, error) {
	if s == "-" {
		return streamid.ID{Time: 0, Seq: 0}, nil
	}
	if s == "+" {
		return streamid.ID{Time: 1<<63 - 1, Seq: 1<<63 - 1}, nil
	}

	if strings.Contains(s, "-") {
		spec, err := ParseSpec(s)
		if err != nil {
			return streamid.ID{}, err
		}
		if spec.Kind != SpecExplicit {
			return streamid.ID{}, rerrors.ErrInvalidArgs
		}
		return streamid.ID{Time: spec.Time, Seq: spec.Seq}, nil
	}

	t, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return streamid.ID{}, rerrors.ErrInvalidArgs
	}

	if isStart {
		return streamid.ID{Time: t, Seq: 0}, nil
	}
	return streamid.ID{Time: t, Seq: 1<<63 - 1}, nil
}

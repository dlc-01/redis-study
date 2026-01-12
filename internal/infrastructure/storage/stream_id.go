package storage

import (
	"fmt"
	"strconv"
	"strings"
)

type streamID struct {
	Time int64
	Seq  int64
}

type idSpecKind int

const (
	specExplicit idSpecKind = iota
	specAutoSeq
	specAutoID
)

type idSpec struct {
	kind idSpecKind
	t    int64
	seq  int64
}

func parseIDSpec(s string) (idSpec, error) {
	if s == "*" {
		return idSpec{kind: specAutoID}, nil
	}

	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return idSpec{}, fmt.Errorf("ERR Invalid stream ID specified as stream command argument")
	}

	t, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return idSpec{}, fmt.Errorf("ERR Invalid stream ID specified as stream command argument")
	}

	if parts[1] == "*" {
		return idSpec{kind: specAutoSeq, t: t}, nil
	}

	seq, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return idSpec{}, fmt.Errorf("ERR Invalid stream ID specified as stream command argument")
	}

	return idSpec{kind: specExplicit, t: t, seq: seq}, nil
}

func (id streamID) String() string {
	return fmt.Sprintf("%d-%d", id.Time, id.Seq)
}

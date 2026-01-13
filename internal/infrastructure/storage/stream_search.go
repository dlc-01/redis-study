package storage

import (
	"github.com/codecrafters-io/redis-starter-go/internal/domain/streamid"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/value"
)

func lowerBoundGE(entries []value.StreamEntry, target streamid.ID) int {
	lo := 0
	hi := len(entries)

	for lo < hi {
		mid := lo + (hi-lo)/2
		if streamid.Compare(entries[mid].ID, target) < 0 {
			lo = mid + 1
		} else {
			hi = mid
		}
	}

	return lo
}

func lowerBoundGT(entries []value.StreamEntry, target streamid.ID) int {
	lo := 0
	hi := len(entries)

	for lo < hi {
		mid := lo + (hi-lo)/2
		if streamid.Compare(entries[mid].ID, target) <= 0 {
			lo = mid + 1
		} else {
			hi = mid
		}
	}

	return lo
}

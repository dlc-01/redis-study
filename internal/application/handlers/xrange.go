package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/internal/application/commands/stream"
	"github.com/codecrafters-io/redis-starter-go/internal/application/errors"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/response"
	"github.com/codecrafters-io/redis-starter-go/internal/ports"
)

func HandleXRange(storage ports.Storage, c stream.XRangeCommand) (response.Response, error) {
	entries, err := storage.XRange(c.Key, c.Start, c.End)
	if err != nil {
		return errors.ToResponse(err), nil
	}

	out := make([]response.Response, 0, len(entries))

	for _, e := range entries {
		id := e.IDRaw

		kvs := make([]response.Response, 0, len(e.Values)*2)
		for _, kv := range e.Values {
			k := kv.Key
			v := kv.Value
			kvs = append(kvs,
				response.BulkString{Value: &k},
				response.BulkString{Value: &v},
			)
		}

		out = append(out, response.Array{
			Values: []response.Response{
				response.BulkString{Value: &id},
				response.Array{Values: kvs},
			},
		})
	}

	return response.Array{Values: out}, nil
}

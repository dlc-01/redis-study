package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/internal/application/commands/stream"
	apperrors "github.com/codecrafters-io/redis-starter-go/internal/application/errors"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/response"
	"github.com/codecrafters-io/redis-starter-go/internal/ports"
)

func HandleXRead(storage ports.Storage, c stream.XReadCommand) (response.Response, error) {
	res, err := storage.XReadMany(c.Keys, c.IDs)
	if err != nil {
		return apperrors.ToResponse(err), nil
	}

	streams := make([]response.Response, 0)

	for _, key := range c.Keys {
		entries, ok := res[key]
		if !ok || len(entries) == 0 {
			continue
		}

		entriesResp := make([]response.Response, 0, len(entries))
		for _, e := range entries {
			id := e.ID

			kvs := make([]response.Response, 0, len(e.Values)*2)
			for _, kv := range e.Values {
				k := kv.Key
				v := kv.Value
				kvs = append(kvs,
					response.BulkString{Value: &k},
					response.BulkString{Value: &v},
				)
			}

			entriesResp = append(entriesResp, response.Array{
				Values: []response.Response{
					response.BulkString{Value: &id},
					response.Array{Values: kvs},
				},
			})
		}

		k := key
		streams = append(streams, response.Array{
			Values: []response.Response{
				response.BulkString{Value: &k},
				response.Array{Values: entriesResp},
			},
		})
	}

	return response.Array{Values: streams}, nil
}

package codec

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/domain/response"
)

type BuilderResponse struct{}

func (b BuilderResponse) Build(resp response.Response) ([]byte, error) {
	switch r := resp.(type) {

	case response.SimpleString:
		return []byte(fmt.Sprintf("+%s\r\n", r.Value)), nil

	case response.BulkString:
		if r.Value == nil {
			return []byte("$-1\r\n"), nil
		}
		s := *r.Value
		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)), nil
	case response.ErrorString:
		return []byte(fmt.Sprintf("-%s\r\n", r.Message)), nil
	case response.Integer:
		return []byte(fmt.Sprintf(":%d\r\n", r.Value)), nil

	}

	return nil, fmt.Errorf("unsupported response type")
}

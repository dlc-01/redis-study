package codec

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/domain/response"
)

type BuilderResponse struct{}

func (b BuilderResponse) Build(res response.Response) ([]byte, error) {
	switch v := res.(type) {

	case response.SimpleString:
		return []byte(fmt.Sprintf("+%s\r\n", v.Value)), nil

	case response.BulkString:
		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v.Value), v.Value)), nil

	default:
		return nil, fmt.Errorf("unknown response type %T", v)
	}
}

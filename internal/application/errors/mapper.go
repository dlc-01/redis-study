package errors

import (
	"errors"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/domain/rerrors"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/response"
)

func ToResponse(err error) response.Response {
	if errors.Is(err, rerrors.ErrWrongType) {
		return response.ErrorString{Message: rerrors.ErrWrongType.Error()}
	}

	if errors.Is(err, rerrors.ErrInvalidArgs) {
		return response.ErrorString{Message: rerrors.ErrInvalidArgs.Error()}
	}

	if errors.Is(err, rerrors.ErrUnknown) {
		return response.ErrorString{Message: rerrors.ErrUnknown.Error()}
	}

	if errors.Is(err, rerrors.ErrInternal) {
		return response.ErrorString{Message: rerrors.ErrInternal.Error()}
	}

	msg := err.Error()
	if strings.HasPrefix(msg, "ERR ") || strings.HasPrefix(msg, "WRONGTYPE ") {
		return response.ErrorString{Message: msg}
	}

	return response.ErrorString{Message: "ERR " + msg}
}

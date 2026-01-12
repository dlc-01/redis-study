package application

import (
	"github.com/codecrafters-io/redis-starter-go/internal/domain/rerrors"
	"github.com/codecrafters-io/redis-starter-go/internal/domain/response"
)

type HandlerFunc func(cmd Command) (response.Response, error)

type Router struct {
	handlers map[string]HandlerFunc
}

func NewRouter() *Router {
	return &Router{handlers: make(map[string]HandlerFunc)}
}

func (r *Router) Register(name string, h HandlerFunc) {
	r.handlers[name] = h
}

func (r *Router) Handle(cmd Command) (response.Response, error) {
	h, ok := r.handlers[cmd.Name()]
	if !ok {
		return nil, rerrors.ErrUnknown
	}
	return h(cmd)
}

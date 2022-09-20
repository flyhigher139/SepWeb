package errhdl

import (
	"github.com/igevin/sepweb/pkg/context"
	"github.com/igevin/sepweb/pkg/handler"
	"github.com/igevin/sepweb/pkg/middleware"
)

type MiddlewareBuilder struct {
	resp map[int][]byte
}

func (m *MiddlewareBuilder) RegisterError(code int, resp []byte) *MiddlewareBuilder {
	m.resp[code] = resp
	return m
}

func (m *MiddlewareBuilder) Build() middleware.Middleware {
	return func(next handler.Handle) handler.Handle {
		return func(ctx *context.Context) {
			next(ctx)
			if resp, ok := m.resp[ctx.RespStatusCode]; ok {
				ctx.RespData = resp
			}
		}
	}
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{resp: make(map[int][]byte, 64)}
}

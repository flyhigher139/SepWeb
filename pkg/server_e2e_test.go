package sepweb

import (
	"net/http"
	"testing"
)

func TestHttpServer_Start(t *testing.T) {
	s := NewHttpServer()
	s.addRoute(http.MethodGet, "/", handleRoot)
	s.addRoute(http.MethodGet, "/static", handleStatic)
	s.addRoute(http.MethodGet, "/home/*", handleStar)
	s.addRoute(http.MethodGet, "/users/:user", handleParam)
	_ = s.Start(":8081")
}

func handleRoot(ctx *Context) {
	ctx.Resp.WriteHeader(http.StatusOK)
	_, _ = ctx.Resp.Write([]byte("hello, root"))
}

func handleStatic(ctx *Context) {
	ctx.Resp.WriteHeader(http.StatusOK)
	_, _ = ctx.Resp.Write([]byte("hello, static"))
}

func handleStar(ctx *Context) {
	ctx.Resp.WriteHeader(http.StatusOK)
	_, _ = ctx.Resp.Write([]byte("hello, home"))
}

func handleParam(ctx *Context) {
	ctx.Resp.WriteHeader(http.StatusOK)
	_, _ = ctx.Resp.Write([]byte("hello, user"))
}

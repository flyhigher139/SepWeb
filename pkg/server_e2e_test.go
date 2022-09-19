package sepweb

import (
	"fmt"
	"github.com/igevin/sepweb/pkg/context"
	"net/http"
	"testing"
)

func TestHttpServer_Start(t *testing.T) {
	s := NewHttpServer()
	s.AddRoute(http.MethodGet, "/", handleRoot)
	s.AddRoute(http.MethodGet, "/static", handleStatic)
	s.AddRoute(http.MethodGet, "/home/*", handleStar)
	s.AddRoute(http.MethodGet, "/users/:user", handleParam)
	_ = s.Start(":8081")
}

func handleRoot(ctx *context.Context) {
	ctx.Resp.WriteHeader(http.StatusOK)
	_, _ = ctx.Resp.Write([]byte("hello, root"))
}

func handleStatic(ctx *context.Context) {
	ctx.Resp.WriteHeader(http.StatusOK)
	_, _ = ctx.Resp.Write([]byte("hello, static"))
}

func handleStar(ctx *context.Context) {
	ctx.Resp.WriteHeader(http.StatusOK)
	_, _ = ctx.Resp.Write([]byte("hello, home"))
}

func handleParam(ctx *context.Context) {
	value := ctx.PathValue("user")
	user, err := value.ToString()
	if err != nil {
		defaultHandleError(ctx)
		return
	}
	ctx.Resp.WriteHeader(http.StatusOK)
	_, _ = ctx.Resp.Write([]byte(fmt.Sprintf("hello, %s", user)))
}

func defaultHandleError(ctx *context.Context) {
	ctx.Resp.WriteHeader(http.StatusBadRequest)
	_, _ = ctx.Resp.Write([]byte("bad request"))
}

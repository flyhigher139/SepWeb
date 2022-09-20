package errhdl

import (
	sepweb "github.com/igevin/sepweb/pkg"
	"github.com/igevin/sepweb/pkg/context"
	"net/http"
	"testing"
)

func TestHttpServer_Start(t *testing.T) {
	s := sepweb.NewHttpServer()
	s.Use(CreateHttpErrorHandleMiddleware())
	s.AddRoute(http.MethodGet, "/home", func(ctx *context.Context) {
		ctx.Resp.WriteHeader(http.StatusOK)
		_, _ = ctx.Resp.Write([]byte("hello, home"))
	})
	_ = s.Start(":8081")
}

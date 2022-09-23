package sepweb

import (
	"fmt"
	"github.com/igevin/sepweb/pkg/context"
	"github.com/igevin/sepweb/pkg/template"
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

func TestServerWithRenderEngine(t *testing.T) {
	engine := &template.GoTemplateEngine{}
	err := engine.LoadFromGlob("testdata/tpls/*.gohtml")
	if err != nil {
		t.Fatal(err)
	}
	s := NewHttpServer(ServerWithTemplateEngine(engine))
	s.Get("/login", func(ctx *context.Context) {
		er := ctx.Render("login.gohtml", nil)
		if er != nil {
			t.Fatal(er)
		}
	})
	err = s.Start(":8081")
	if err != nil {
		t.Fatal(err)
	}
}

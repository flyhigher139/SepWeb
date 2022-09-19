package sepweb

import (
	"github.com/igevin/sepweb/pkg/context"
	"github.com/igevin/sepweb/pkg/route"
	"net/http"
)

type Server interface {
	http.Handler
	Start(addr string) error
}

type HttpServer struct {
	route.Router
}

func (s *HttpServer) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := &context.Context{
		Req:  r,
		Resp: w,
	}
	s.serve(ctx)
}

func (s *HttpServer) serve(ctx *context.Context) {
	mi, ok := s.FindRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || mi == nil {
		ctx.Resp.WriteHeader(http.StatusNotFound)
		_, _ = ctx.Resp.Write([]byte("Not Found"))
		return
	}
	ctx.PathParams = mi.PathParams
	mi.N.Handler(ctx)
}

func NewHttpServer() *HttpServer {
	return &HttpServer{Router: route.NewRouter()}
}

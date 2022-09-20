package sepweb

import (
	"github.com/igevin/sepweb/pkg/context"
	"github.com/igevin/sepweb/pkg/handler"
	"github.com/igevin/sepweb/pkg/middleware"
	"github.com/igevin/sepweb/pkg/route"
	"log"
	"net/http"
)

type Server interface {
	http.Handler
	Start(addr string) error
}

type HttpServer struct {
	route.Router
	mdls []middleware.Middleware
}

func (s *HttpServer) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := &context.Context{
		Req:  r,
		Resp: w,
	}
	handle := s.serve
	for i := len(s.mdls) - 1; i >= 0; i-- {
		handle = s.mdls[i](handle)
	}

	handle = s.flushRespMiddleware(handle)
	handle(ctx)
}

func (s *HttpServer) Use(mdls ...middleware.Middleware) {
	if s.mdls == nil {
		s.mdls = mdls
		return
	}
	s.mdls = append(s.mdls, mdls...)
}

func (s *HttpServer) serve(ctx *context.Context) {
	mi, ok := s.FindRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || mi.N == nil || mi.N.Handler == nil {
		ctx.RespStatusCode = http.StatusNotFound
		ctx.RespData = []byte("Not Found")
		return
	}
	ctx.PathParams = mi.PathParams
	ctx.MatchedRoute = mi.N.Route
	mi.N.Handler(ctx)
}

func (s *HttpServer) flushResp(ctx *context.Context) {
	if ctx.RespStatusCode > 0 {
		ctx.Resp.WriteHeader(ctx.RespStatusCode)
	}
	_, err := ctx.Resp.Write(ctx.RespData)
	if err != nil {
		log.Fatalln("回写响应失败", err)
	}
}

func (s *HttpServer) flushRespMiddleware(next handler.Handle) handler.Handle {
	return func(ctx *context.Context) {
		next(ctx)
		s.flushResp(ctx)
	}
}

func NewHttpServer() *HttpServer {
	return &HttpServer{Router: route.NewRouter()}
}

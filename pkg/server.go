package sepweb

import (
	"github.com/igevin/sepweb/pkg/context"
	"github.com/igevin/sepweb/pkg/handler"
	"github.com/igevin/sepweb/pkg/middleware"
	"github.com/igevin/sepweb/pkg/route"
	"github.com/igevin/sepweb/pkg/template"
	"log"
	"net/http"
)

type Server interface {
	http.Handler
	Start(addr string) error
}

type HttpServer struct {
	route.Router
	mdls      []middleware.Middleware
	tplEngine template.TemplateEngine
}

type ServerOption func(server *HttpServer)

func (s *HttpServer) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := &context.Context{
		Req:       r,
		Resp:      w,
		TplEngine: s.tplEngine,
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

func (s *HttpServer) Get(path string, handle handler.Handle) {
	s.AddRoute(http.MethodGet, path, handle)
}

func (s *HttpServer) Post(path string, handle handler.Handle) {
	s.AddRoute(http.MethodPost, path, handle)
}

func NewHttpServer(opts ...ServerOption) *HttpServer {
	s := &HttpServer{Router: route.NewRouter()}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func ServerWithTemplateEngine(engine template.TemplateEngine) ServerOption {
	return func(server *HttpServer) {
		server.tplEngine = engine
	}
}

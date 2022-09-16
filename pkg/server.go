package sepweb

import "net/http"

type HandlerFunc func(ctx *Context)

type Server interface {
	http.Handler
	Start(addr string) error
	//AddRoute(method string, path string, handler HandlerFunc)
}

type HttpServer struct {
	router
}

func (s *HttpServer) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := &Context{
		Req:  r,
		Resp: w,
	}
	s.serve(ctx)
}

func (s *HttpServer) serve(ctx *Context) {
	matchInfo, ok := s.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || matchInfo == nil {
		ctx.Resp.WriteHeader(http.StatusNotFound)
		_, _ = ctx.Resp.Write([]byte("Not Found"))
		return
	}
	matchInfo.n.handler(ctx)
}

func NewHttpServer() *HttpServer {
	return &HttpServer{router: newRouter()}
}

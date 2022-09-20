package route

import (
	"fmt"
	"github.com/igevin/sepweb/pkg/handler"
	"strings"
)

type Router struct {
	routes map[string]*node
}

func NewRouter() Router {
	return Router{
		routes: map[string]*node{},
	}
}

func (r *Router) AddRoute(method, path string, handler handler.Handle) {
	_ = r.checkPathFormat(path)
	root, ok := r.handleRootRouter(method, path, handler)
	if ok {
		return
	}
	r.handleSegmentRouter(root, path, handler)
}

func (r *Router) checkPathFormat(path string) bool {
	if path == "" {
		panic("web: 路由是空字符串")
	}
	if path[0] != '/' {
		panic("web: 路由必须以 / 开头")
	}

	if path != "/" && path[len(path)-1] == '/' {
		panic("web: 路由不能以 / 结尾")
	}
	return true
}

func (r *Router) handleRootRouter(method, path string, handler handler.Handle) (*node, bool) {
	root, ok := r.routes[method]
	if !ok {
		root = &node{path: "/"}
		r.routes[method] = root
	}
	if path == "/" {
		if root.Handler != nil {
			panic("web: 路由冲突[/]")
		}
		root.Handler = handler
		root.Route = path
	}
	return root, path == "/"
}

func (r *Router) handleSegmentRouter(n *node, path string, handler handler.Handle) {
	segs := strings.Split(path[1:], "/")
	// 开始一段段处理
	for _, s := range segs {
		if s == "" {
			panic(fmt.Sprintf("web: 非法路由。不允许使用 //a/b, /a//b 之类的路由, [%s]", path))
		}
		n = n.childOrCreate(s)
	}
	if n.Handler != nil && n.path != "*" {
		panic(fmt.Sprintf("web: 路由冲突[%s]", path))
	}
	n.Handler = handler
	n.Route = path
}

func (r *Router) FindRoute(method, path string) (*matchInfo, bool) {
	root, ok := r.routes[method]
	if !ok {
		return &matchInfo{}, false
	}

	if path == "/" {
		return &matchInfo{N: root}, true
	}

	return r.findPathRoute(root, path)
}

func (r *Router) findPathRoute(curNode *node, path string) (*matchInfo, bool) {
	segs := strings.Split(strings.Trim(path, "/"), "/")
	mi := &matchInfo{}
	var prev *node
	for _, s := range segs {
		var matchParam, ok bool
		prev = curNode
		curNode, matchParam, ok = curNode.childOf(s)
		if !ok && prev.typ != nodeTypeAny {
			return nil, false
		}
		if curNode == nil {
			curNode = prev
		}
		if matchParam {
			mi.addValue(curNode.paramName, s)
		}
	}
	mi.N = curNode
	return mi, true
}

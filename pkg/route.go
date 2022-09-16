package pkg

import (
	"fmt"
	"regexp"
	"strings"
)

type router struct {
	routes map[string]*node
}

func newRouter() router {
	return router{
		routes: map[string]*node{},
	}
}

func (r *router) addRoute(method, path string, handler HandlerFunc) {
	_ = r.checkPathFormat(path)
	root, ok := r.handleRootRouter(method, path, handler)
	if ok {
		return
	}
	r.handleSegmentRouter(root, path, handler)
}

func (r *router) checkPathFormat(path string) bool {
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

func (r *router) handleRootRouter(method, path string, handler HandlerFunc) (*node, bool) {
	root, ok := r.routes[method]
	if !ok {
		root = &node{path: "/"}
		r.routes[method] = root
	}
	if path == "/" {
		if root.handler != nil {
			panic("web: 路由冲突[/]")
		}
		root.handler = handler
	}
	return root, path == "/"
}

func (r *router) handleSegmentRouter(n *node, path string, handler HandlerFunc) {
	segs := strings.Split(path[1:], "/")
	// 开始一段段处理
	for _, s := range segs {
		if s == "" {
			panic(fmt.Sprintf("web: 非法路由。不允许使用 //a/b, /a//b 之类的路由, [%s]", path))
		}
		n = n.childOrCreate(s)
	}
	if n.handler != nil && n.path != "*" {
		panic(fmt.Sprintf("web: 路由冲突[%s]", path))
	}
	n.handler = handler
}

func (r *router) findRoute(method, path string) (*matchInfo, bool) {
	root, ok := r.routes[method]
	if !ok {
		return nil, false
	}

	if path == "/" {
		return &matchInfo{n: root}, true
	}

	return r.findPathRoute(root, path)
}

func (r *router) findPathRoute(curNode *node, path string) (*matchInfo, bool) {
	segs := strings.Split(strings.Trim(path, "/"), "/")
	mi := &matchInfo{}
	var prev *node
	for _, s := range segs {
		var matchParam, ok bool
		prev = curNode
		curNode, matchParam, ok = curNode.childOf(s)
		if !ok {
			return nil, false
		}
		if curNode == nil {
			curNode = prev
		}
		if matchParam {
			mi.addValue(curNode.paramName, s)
		}
	}
	mi.n = curNode
	return mi, true
}

type nodeType int

const (
	// 静态路由
	nodeTypeStatic = iota
	// 正则路由
	nodeTypeReg
	// 路径参数路由
	nodeTypeParam
	// 通配符路由
	nodeTypeAny
)

type node struct {
	path     string
	children map[string]*node
	typ      nodeType

	starChild *node

	paramChild *node
	paramName  string

	regChild *node
	regExpr  *regexp.Regexp

	handler HandlerFunc
}

func (n *node) childOrCreate(path string) *node {
	if path == "*" {
		return n.starChildOrCreate(path)
	}

	paramName, regExpr := n.matchAndParseRegExp(path)
	// 解析到正则，是正则路由
	if regExpr != nil {
		return n.regexChildOrCreate(path, paramName, regExpr)
	}

	// 以 : 开头，我们认为是参数路由
	if path[0] == ':' {
		return n.paramChildOrCreate(path)
	}

	return n.staticChildOrCreate(path)
}

func (n *node) starChildOrCreate(path string) *node {
	_ = n.isStarChildAvailable(path)
	if n.starChild == nil {
		n.starChild = &node{path: path, typ: nodeTypeAny}
	}
	return n.starChild
}

func (n *node) isStarChildAvailable(path string) bool {
	if n.paramChild != nil {
		panic(fmt.Sprintf("web: 非法路由，已有路径参数路由。不允许同时注册通配符路由和参数路由 [%s]", path))
	}
	if n.regChild != nil {
		panic(fmt.Sprintf("web: 非法路由，已有正则路由。不允许同时注册通配符路由和正则路由 [%s]", path))
	}
	return true
}

func (n *node) matchAndParseRegExp(path string) (string, *regexp.Regexp) {
	if !strings.HasPrefix(path, ":") || !strings.HasSuffix(path, ")") || !strings.Contains(path, "(") {
		return "", nil
	}
	segs := strings.SplitN(path[1:len(path)-1], "(", 2)

	if reg := regexp.MustCompile(segs[1]); reg == nil {
		return "", nil
	} else {
		return segs[0], reg
	}
}

func (n *node) regexChildOrCreate(path, paramName string, regExpr *regexp.Regexp) *node {
	_ = n.isRegexChildAvailable(path)
	if n.regChild == nil {
		n.regChild = &node{path: path, paramName: paramName, regExpr: regExpr, typ: nodeTypeReg}
	}

	return n.regChild
}

func (n *node) isRegexChildAvailable(path string) bool {
	if n.starChild != nil {
		panic(fmt.Sprintf("web: 非法路由，已有通配符路由。不允许同时注册通配符路由和正则路由 [%s]", path))
	}
	if n.paramChild != nil {
		panic(fmt.Sprintf("web: 非法路由，已有路径参数路由。不允许同时注册正则路由和参数路由 [%s]", path))
	}
	if n.regChild != nil && n.regChild.path != path {
		panic(fmt.Sprintf("web: 路由冲突，参数路由冲突，已有 %s，新注册 %s", n.regChild.path, path))
	}

	return true
}

func (n *node) paramChildOrCreate(path string) *node {
	_ = n.isParamChildAvailable(path)
	if n.paramChild == nil {
		n.paramChild = &node{path: path, paramName: path[1:], typ: nodeTypeParam}
	}
	return n.paramChild
}

func (n *node) isParamChildAvailable(path string) bool {
	if n.starChild != nil {
		panic(fmt.Sprintf("web: 非法路由，已有通配符路由。不允许同时注册通配符路由和参数路由 [%s]", path))
	}
	if n.regChild != nil {
		panic(fmt.Sprintf("web: 非法路由，已有正则路由。不允许同时注册正则路由和参数路由 [%s]", path))
	}
	if n.paramChild != nil && n.paramChild.path != path {
		panic(fmt.Sprintf("web: 路由冲突，参数路由冲突，已有 %s，新注册 %s", n.paramChild.path, path))
	}

	return true
}

func (n *node) staticChildOrCreate(path string) *node {
	if n.children == nil {
		n.children = make(map[string]*node)
	}
	child, ok := n.children[path]
	if !ok {
		child = &node{path: path, typ: nodeTypeStatic}
		n.children[path] = child
	}
	return child
}

func (n *node) childOf(path string) (*node, bool, bool) {
	if n.children == nil {
		if n.regChild != nil {
			matched := n.regChild.regExpr.MatchString(path)
			return n.regChild, matched, true
		}
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, true
		//return n.starChild, false, n.starChild != nil
	}
	res, ok := n.children[path]
	if !ok {
		if n.regChild != nil {
			matched := n.regChild.regExpr.MatchString(path)
			return n.regChild, true, matched
		}
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}
	return res, false, ok
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}

func (m *matchInfo) addValue(key string, value string) {
	if m.pathParams == nil {
		m.pathParams = make(map[string]string, 1)
	}
	m.pathParams[key] = value
}

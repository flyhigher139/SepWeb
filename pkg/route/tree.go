package route

import (
	"fmt"
	"github.com/igevin/sepweb/pkg/handler"
	"regexp"
	"strings"
)

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

	Handler handler.Handle
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
		//return N.starChild, false, N.starChild != nil
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

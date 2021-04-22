package gen

import (
	"fmt"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandleFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandleFunc),
	}
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' { // 静态文件直接退出
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandleFunc) {
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}

	parts := parsePattern(pattern)
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[routeKey(method, pattern)] = handler
}

func routeKey(method string, pattern string) string {
	return method + "-" + pattern
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	// 分析出part
	searchParts := parsePattern(path)
	params := make(map[string]string)

	n := root.search(searchParts, 0)
	if n != nil { // 说明找到了匹配的路由，获取param的值
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index] // 动态路由的param保存
			}

			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/") //
				break
			}
		}
	}
	return n, params
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	fmt.Printf("%+v, %+v\n\n", n, params)
	if n != nil {
		c.Params = params
		r.handlers[routeKey(c.Method, n.pattern)](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}

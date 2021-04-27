package gen

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 分析路由
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

// 添加路由
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}

	parts := parsePattern(pattern)
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[routeKey(method, pattern)] = handler
}

// 获取map中的路由key
func routeKey(method string, pattern string) string {
	return method + "-" + pattern
}

// 根据请求方法+路径获取路由，并解析路由中的参数绑定
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

			if part[0] == '*' && len(part) > 1 { // 如果这个node的part[0]是*，则将后续path作为一个path返回
				params[part[1:]] = strings.Join(searchParts[index:], "/") //
				break
			}
		}
	}
	return n, params
}

func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

// 用来处理请求
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path) // 解析路由
	if n != nil {                             // 路由存在，则执行对应的处理逻辑
		c.Params = params
		c.handlers = append(c.handlers, r.handlers[routeKey(c.Method, n.pattern)])
	} else { // 路由不存在，则404
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}

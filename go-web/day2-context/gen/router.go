package gen

import "net/http"

type router struct {
	handlers map[string]HandleFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandleFunc)}
}

func (r *router) addRoute(method string, pattern string, handler HandleFunc) {
	r.handlers[routeKey(method, pattern)] = handler
}

func routeKey(method string, pattern string) string {
	return method + "-" + pattern
}

func (r *router) handle(c *Context) {
	if handler, ok := r.handlers[routeKey(c.Method, c.Path)]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}

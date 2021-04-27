package gen

import (
	"log"
	"net/http"
	"strings"
)

/**
框架基类，贯穿整个应用，包括服务启动，端口监听，路由配置，路由匹配，中间件注册等一系列工作
*/
type HandlerFunc func(*Context)

type (
	// 根据实际情况拆分，中间件是以group为维度，故放在这里
	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc
		parent      *RouterGroup
		engine      *Engine
	}

	// 贯穿整个应用的结构体
	Engine struct {
		*RouterGroup // 相当于继承RouterGroup
		router       *router
		groups       []*RouterGroup
	}
)

// 实例Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// 实例化带Logger中间件的Engine
func Default() *Engine {
	engine := New()
	engine.Use(Logger())
	return engine
}

// 获取新的RouterGroup
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 基于RouterGroup添加路由
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Router %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// 设置GET类路由
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// 设置POST路由
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// 代理http，执行监听
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// 监听到请求时，执行的回调
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		// 这里直接要卖URL.Path是否包含group.prefix来判断是否使用这个group的中间件
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares // 初始化中间件
	engine.router.handle(c)
}

package gee

import "net/http"

type RouterGroup struct {
	prefix     string
	middleware []HandleFunc
	engine     *Engine
}

func (routerGroup *RouterGroup) GET(pattern string, handle HandleFunc) {
	routerGroup.addRouter(http.MethodGet, pattern, handle)
}
func (routerGroup *RouterGroup) POST(pattern string, handle HandleFunc) {
	routerGroup.addRouter(http.MethodPost, pattern, handle)
}
func (routerGroup *RouterGroup) addRouter(method string, pattern string, handle HandleFunc) {
	routerGroup.engine.routers.addRouter(method, routerGroup.prefix+pattern, handle)
}

func (routerGroup *RouterGroup) Group(prefix string) *RouterGroup {
	engine := routerGroup.engine

	group := new(RouterGroup)
	group.prefix += prefix
	group.engine = engine

	engine.groups = append(engine.groups, group)
	return group
}

func (routerGroup *RouterGroup) Use(middleware ...HandleFunc) {
	if routerGroup.middleware == nil {
		routerGroup.middleware = make([]HandleFunc, 0)
	}

	routerGroup.middleware = append(routerGroup.middleware, middleware...)
}

func (routerGroup *RouterGroup) Static(pattern, root string) {
	// 创建静态文件请求处理器
	handler := routerGroup.createStaticHandler(pattern, root)
	routerGroup.GET(pattern+"/*filepath", handler)
}

// createStaticHandler 创建静态文件处理器
// @pattern 请求目录的prefix
// @fileSystem 资源文件的真是目录
func (routerGroup *RouterGroup) createStaticHandler(pattern, root string) HandleFunc {
	reqPath := routerGroup.prefix + pattern
	fileSystem := http.Dir(root)
	// StripPrefix 用于将请求重定向到fileServer之前将 req path 的 prefix 拿掉
	fileServer := http.StripPrefix(reqPath, http.FileServer(fileSystem))

	// 处理器逻辑：从req path 中获取除prefix之外的文件名称并查看文件是否存在
	// 如果不存在则直接返回错误，如果存在则将请求交给httpServer静态服务器
	return func(c *Context) {
		fileName, ok := c.ParamGet("filepath")
		if !ok {
			c.String(http.StatusNotFound, "Not found file:"+fileName+"("+reqPath+" -> "+root+")")
			return
		}

		if _, err := fileSystem.Open(fileName); err != nil {
			c.String(http.StatusNotFound, "Not found file:"+fileName+"("+reqPath+" -> "+root+")")
			return
		}

		fileServer.ServeHTTP(c.writer, c.req)
	}
}

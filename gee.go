// gee 框架核心实现
package gee

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"
)

type HandleFunc func(c *Context)

type (
	Engine struct {
		*RouterGroup
		routers *Router
		groups  []*RouterGroup
		// 定义了将所有的模板都 load 进 engine 对象
		templates *template.Template

		// 自定义模板渲染函数的列表
		funcMap template.FuncMap
	}
	RouterGroup struct {
		prefix     string
		middleware []HandleFunc
		engine     *Engine
	}
)

func New() *Engine {
	// 循环指针: Engine <-> RouterGroup
	engine := new(Engine)
	engine.routers = NewRouter()
	engine.groups = make([]*RouterGroup, 0)

	group := new(RouterGroup)
	group.engine = engine

	engine.RouterGroup = group
	engine.groups = append(engine.groups, group)
	return engine
}

func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	//engine.Use(Recovery())
	return engine
}

func (engine *Engine) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// 将当前请求匹配的middleware加入context对象中
	middlewares := make([]HandleFunc, 0)
	for _, group := range engine.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middleware...)
		}
	}

	c := NewContext(rw, r, middlewares, engine)
	engine.routers.handle(c)
}
func (engine *Engine) Run(addr string) error {
	log.Println("[Gee Web] Server listening at " + addr)
	return http.ListenAndServe(addr, engine)
}
func (engine *Engine) SetFuncMap(fm template.FuncMap) {
	engine.funcMap = fm
}

// LoadHTMLGlob 将所有的模板都加载进engine对象中
//
func (engine *Engine) LoadHTMLGlob(pattern string) {
	glob, err := template.New("").Funcs(engine.funcMap).ParseGlob(pattern)
	engine.templates = template.Must(glob, err)
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

		fileServer.ServeHTTP(c.Rw, c.Req)
	}
}

func Logger() HandleFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Printf("[Gee Web] Logger for all: [%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func Recovery() HandleFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				var logPrint strings.Builder
				logPrint.WriteString(fmt.Sprintf("[Gee Web] [Panic] %s trace back:\n", err))
				var pcs [32]uintptr
				// 跳过前3个caller
				n := runtime.Callers(3, pcs[:])
				for _, pc := range pcs[:n] {
					// 获取对应函数
					fun := runtime.FuncForPC(pc)
					// 获取函数文件名和行号
					fileName, line := fun.FileLine(pc)
					logPrint.WriteString(fmt.Sprintf("\t%s:%d\n", fileName, line))
				}
				log.Println(logPrint.String())
				c.String(http.StatusInternalServerError, "Internal server error")
			}
		}()
		c.Next()
	}
}

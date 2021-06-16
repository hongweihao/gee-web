// gee 框架核心实现
package gee

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

type (
	HandleFunc func(c *Context)
	H          map[string]interface{}
)

type Engine struct {
	*RouterGroup
	routers *Router
	groups  []*RouterGroup
	// 定义了将所有的模板都 load 进 engine 对象
	templates *template.Template

	// 自定义模板渲染函数的列表
	funcMap template.FuncMap
}

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

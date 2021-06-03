// 上下文 Context
package gee

import (
	"encoding/json"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Req        *http.Request
	Rw         http.ResponseWriter
	Path       string
	Method     string
	Param      map[string]string
	StatusCode int

	// middleware
	middlewares      []HandleFunc
	middlewaresIndex int

	engine *Engine
}

func NewContext(rw http.ResponseWriter, req *http.Request, middlewares []HandleFunc, engine *Engine) *Context {
	return &Context{
		Req:              req,
		Rw:               rw,
		Path:             req.URL.Path,
		Method:           req.Method,
		middlewares:      middlewares,
		middlewaresIndex: -1,
		engine:           engine,
	}
}

func (c *Context) PostForm(key string) string {
	c.Req.ParseForm()
	return c.Req.PostForm.Get(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
}

func (c *Context) SetHeader(key string, value string) {
	c.Rw.Header().Set(key, value)
}

func (c *Context) JSON(code int, rsp interface{}) {
	c.SetHeader("Content-Type", "application/json")
	bs, err := json.Marshal(rsp)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Rw.Write([]byte(err.Error()))
		return
	}
	c.Status(code)
	c.Rw.Write(bs)
}
func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.templates.ExecuteTemplate(c.Rw, name, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
func (c *Context) String(code int, text string) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Rw.Write([]byte(text))
}
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Rw.Write(data)
}
func (c *Context) ParamGet(key string) (string, bool) {
	v, ok := c.Param[key]
	return v, ok
}
func (c *Context) Next() {
	c.middlewaresIndex++
	if c.middlewaresIndex >= len(c.middlewares) {
		return
	}
	for ; c.middlewaresIndex < len(c.middlewares); c.middlewaresIndex++ {
		c.middlewares[c.middlewaresIndex](c)
	}
}

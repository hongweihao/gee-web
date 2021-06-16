// 上下文 Context
package gee

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Request        *http.Request
	responseWriter http.ResponseWriter
	Path           string
	Method         string
	param          map[string]string
	StatusCode     int

	// middleware
	middlewares      []HandleFunc
	middlewaresIndex int

	engine *Engine
}

func NewContext(rw http.ResponseWriter, req *http.Request, middlewares []HandleFunc, engine *Engine) *Context {
	return &Context{
		Request:          req,
		responseWriter:   rw,
		Path:             req.URL.Path,
		Method:           req.Method,
		middlewares:      middlewares,
		middlewaresIndex: -1,
		engine:           engine,
	}
}

func (c *Context) PostForm(key string) string {
	return c.Request.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
}

func (c *Context) SetHeader(key string, value string) {
	c.responseWriter.Header().Set(key, value)
}

func (c *Context) JSON(code int, resp interface{}) {
	c.SetHeader("Content-Type", "application/json")
	bs, err := json.Marshal(resp)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.responseWriter.Write([]byte(err.Error()))
		return
	}
	c.Status(code)
	c.responseWriter.Write(bs)
}
func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.templates.ExecuteTemplate(c.responseWriter, name, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
func (c *Context) String(code int, text string) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.responseWriter.Write([]byte(text))
}
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.responseWriter.Write(data)
}
func (c *Context) Param(key string) string {
	return c.param[key]
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

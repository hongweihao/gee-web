package gee

import (
	"log"
	"net/http"
	"strings"
)

type Router struct {
	roots    map[string]*Trie
	handlers map[string]HandleFunc
}

func NewRouter() *Router {
	return &Router{
		roots:    make(map[string]*Trie),
		handlers: make(map[string]HandleFunc),
	}
}

func (router *Router) addRouter(method, pattern string, handle HandleFunc) {
	tree := router.roots[method]
	if tree == nil {
		root := NewTrie()
		router.roots[method] = root
		tree = root
	}
	tree.Insert(pattern)

	key := strings.Join([]string{method, pattern}, "-")
	router.handlers[key] = handle
	log.Printf("[Gee Web]\t[%s]\t-\t%s\n", method, pattern)
}

func (router *Router) handle(c *Context) {

	tree := router.roots[c.Method]
	n, param := tree.Search(c.Path)
	if n == nil {
		c.middlewares = append(c.middlewares, func(c *Context) {
			c.Status(http.StatusNotFound)
			c.String(http.StatusNotFound, "unknown Path:"+c.Path)
		})
		c.Next()
		return
	}

	key := strings.Join([]string{c.Method, n.pattern}, "-")
	c.param = param
	handler, ok := router.handlers[key]
	if !ok {
		c.middlewares = append(c.middlewares, func(c *Context) {
			c.Status(http.StatusNotFound)
			c.String(http.StatusNotFound, "unknown Path:"+c.Path)
		})
		c.Next()
		return
	}

	c.middlewares = append(c.middlewares, handler)
	c.Next()
}

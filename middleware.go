package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"
)

func Logger() HandleFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Printf("[Gee Web] [%d] %s in %v", c.statusCode, c.req.RequestURI, time.Since(t))
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

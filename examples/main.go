package main

import (
	"fmt"
	"github.com/hongweihao/gee-web"
	"html/template"
	"io/ioutil"
	"log"
	"time"
)

func main() {

	// new engine 对象
	//gee.New()

	// new engine 对象并默认使用logger和recovery中间件
	engine := gee.Default()

	// 使用中间件
	engine.Use(printRequest)

	// 使用GET，POST
	engine.GET("/test", test)
	engine.GET("/test/:id", testId)
	engine.POST("/test_post", testPost)

	// 使用 Group
	userGroup := engine.Group("/user")
	userGroup.GET("/list", UserList)
	userGroup.POST("/create", UserCreate)

	// 使用静态文件。注意设置 working directory 到 main 文件所在目录，否则会出现404
	engine.Static("/assets", "./static")

	// 使用模板
	engine.GET("/index", tmplIndex)
	engine.GET("/index/list", tmplList)

	// 自定义模板函数
	engine.GET("/index/today", tmplToday)

	// 自定义模板处理函数
	engine.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	// 加载模板
	engine.LoadHTMLGlob("templates/*")

	// 启动监听
	err := engine.Run(":8088")
	if err != nil {
		log.Fatal(err)
	}
}

func tmplToday(c *gee.Context) {
	c.HTML(200, "today.html", gee.H{
		"now": time.Now(),
	})
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

type User struct {
	Name string
	Age  uint8
}

func tmplList(c *gee.Context) {
	list := []*User{
		{Name: "mkii", Age: 18},
		{Name: "Karl", Age: 29},
	}
	c.HTML(200, "list.html", gee.H{"list": list})
}

func tmplIndex(c *gee.Context) {
	c.HTML(200, "index.html", nil)
}

func testId(c *gee.Context) {
	id := c.Param("id")
	c.JSON(200, gee.H{
		"id": id,
	})
}

func UserCreate(c *gee.Context) {
	bs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(500, "500")
	}
	c.JSON(200, gee.H{
		"data": string(bs),
	})
}

func UserList(c *gee.Context) {
	c.JSON(200, gee.H{
		"user1": "user_id1",
		"user2": "user_id2",
		"user3": "user_id3",
		"user4": "user_id4",
	})
}

func testPost(c *gee.Context) {
	c.String(200, "test post")
}

func test(c *gee.Context) {
	c.String(200, "test")
}

func printRequest(c *gee.Context) {
	log.Printf("[Gee Web] [%s] %s - %d", c.Request.Method, c.Path, c.Request.ContentLength)
	c.Next()
}

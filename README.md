# gee-web



## 1. 如何安装

```shell
go get github.com/hongweihao/gee-web
```



## 2. 如何使用

### 2.1  启动一个简单的服务

```go
package main
import gee "github.com/hongweihao/gee-web"
func main(){
    engine := gee.Default()
    engine.GET("/test", func (c *gee.Context) {
		c.String(200, "test")
	})
    err := engine.Run(":8088")
	if err != nil {
		log.Fatal(err)
	}
}
// curl http://127.0.0.1:8088/test
```



### 2.2 使用 path 参数

```go
engine.GET("/test/:id", func (c *gee.Context) {
    id := c.Param("id")
    c.JSON(200, gee.H{
        "id": id,
    })
})
// curl http://127.0.0.1:8088/test/1024
```



### 2.3 使用 POST

```go
engine.POST("/test_post", func (c *gee.Context) {
    c.String(200, "post test")
})
// curl http://127.0.0.1:8088/test_post -X POST
```



### 2.4 使用 Group

```go
userGroup := engine.Group("/user")
userGroup.GET("/list", UserList)
userGroup.POST("/create", UserCreate)

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

// curl http://127.0.0.1:8088/user/list
// curl http://127.0.0.1:8088/user/create -X POST -H 'Content-Type:application/json' -d '{"userId": 1024, "name": "mkii"}'
```



### 2.5 使用中间件

```go
engine.Use(printRequest)

func printRequest(c *gee.Context) {
	log.Printf("[Gee Web] [%s] %s - %d", c.Request.Method, c.Path, c.Request.ContentLength)
	c.Next()
}
```



### 2.6 发布静态资源

```go
// assets 的请求会转到去 static 目录下找静态文件
engine.Static("/assets", "./static")
```

> working directory 设置到 project 目录，否则会出现 404



### 2.7 使用页面模板

```go
// 加载扫描模板
engine.LoadHTMLGlob("templates/*")

engine.GET("/index", tmplIndex)
engine.GET("/index/list", tmplList)

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

// curl http://127.0.0.1:8088/index
// curl http://127.0.0.1:8088/index/list
```



templates/index.html：

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Index</title>
    <link rel="stylesheet" href="/assets/style.css"/>
</head>
<body>
    <p>Hello world</p>
</body>
</html>
```

templates/list.html：

```html
<!DOCTYPE html>
<html lang="en">
<title>List</title>
<body>
{{range $index, $ele := .list }}
<p>{{ $index }}: {{ $ele.Name }} is {{ $ele.Age }} years old</p>
{{ end }}
</body>
</html>
```



### 2.8 模板中使用自定义函数

```go
engine.GET("/index/today", tmplToday)

// 自定义模板处理函数
engine.SetFuncMap(template.FuncMap{
    "FormatAsDate": FormatAsDate,
})
// 加载模板
engine.LoadHTMLGlob("templates/*")

func tmplToday(c *gee.Context) {
	c.HTML(200, "today.html", gee.H{
		"now": time.Now(),
	})
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

// curl http://127.0.0.1:8088/index/today
```



templates/today.html：

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Today</title>
</head>
<body>

{{.now | FormatAsDate}}

</body>
</html>
```



## 3. Fork from

[geektutu/7days-golang: 7 days golang programs from scratch (web framework Gee, distributed cache GeeCache, object relational mapping ORM framework GeeORM, rpc framework GeeRPC etc) 7天用Go动手写/从零实现系列 (github.com)](https://github.com/geektutu/7days-golang)

[7天用Go从零实现Web框架Gee教程 | 极客兔兔 (geektutu.com)](https://geektutu.com/post/gee.html)
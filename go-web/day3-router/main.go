package main

// $ curl localhost:8080
// <h1>Hello Gee</h1>
// $ curl localhost:8080/hello?name=lovecucu
// hello lovecucu, you're at /hello
// $ curl localhost:8080/hello/lovecucu
// hello lovecucu, you're at /hello/lovecucu
// $ curl "http://localhost:8080/login" -X POST -d 'username=geektutu&password=1234'
// {"password":"1234","username":"geektutu"}
// $ curl localhost:8080/xxx
// 404 NOT FOUND: /xxx
// curl localhost:8080/assets/a.js
// {"filepath":"a.js"}

import (
	"gen"
	"net/http"
)

func main() {
	engine := gen.New()
	engine.GET("/", func(c *gen.Context) {
		// expect /
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>\n")
	})
	engine.GET("/hello", func(c *gen.Context) {
		// expect /hello?name=lovecucu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})
	engine.GET("/hello/:name", func(c *gen.Context) {
		// expect /hello/lovecucu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})
	engine.GET("/assets/*filepath", func(c *gen.Context) {
		// expect /assets/css/a.css
		c.JSON(http.StatusOK, gen.H{"filepath": c.Param("filepath")})
	})
	engine.POST("/login", func(c *gen.Context) {
		c.JSON(http.StatusOK, gen.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	engine.Run(":8080")
}

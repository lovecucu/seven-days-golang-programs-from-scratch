package main

/*
$ curl -i http://localhost:8080/index
<h1>Index Page</h1>

$ curl -i http://localhost:8080/v1/
<h1>Hello V1 Page</h1>

$ curl http://localhost:8080/v1/hello?name=lovecucu
hello lovecucu, you're at /v1/hello

$ curl http://localhost:8080/v2/hello/lovecucu
hello lovecucu, you're at /v2/hello/lovecucu

$ curl "http://localhost:8080/v2/login" -X POST -d 'username=lovecucu&password=1234'
{"password":"1234","username":"lovecucu"}

$ curl http://localhost:8080/assets/a.cs
{"filepath":"a.cs"}

$ curl http://localhost:8080/hello
404 NOT FOUND: /hello
*/

import (
	"gen"
	"net/http"
)

func main() {
	engine := gen.New()

	engine.GET("/index", func(c *gen.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>\n")
	})

	engine.GET("/assets/*filepath", func(c *gen.Context) {
		// expect /assets/css/a.css
		c.JSON(http.StatusOK, gen.H{"filepath": c.Param("filepath")})
	})

	v1 := engine.Group("/v1")
	{
		v1.GET("/", func(c *gen.Context) {
			c.HTML(http.StatusOK, "<h1>Hello V1 Page</h1>\n")
		})
		v1.GET("/hello", func(c *gen.Context) {
			// expect /hello?name=lovecucu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := engine.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *gen.Context) {
			// expect /hello/lovecucu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})

		v2.POST("/login", func(c *gen.Context) {
			c.JSON(http.StatusOK, gen.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}

	engine.Run(":8080")
}

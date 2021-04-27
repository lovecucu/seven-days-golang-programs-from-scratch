package main

/*
$ curl http://localhost:8080/
<h1>Hello Gen</h1>

server的log:
2021/04/27 19:19:56 [200] / in 2.371µs

$ curl -i http://localhost:8080/v2/lovecucu
{"message":"Internal Server Error"}

server的log:
2021/04/27 19:20:28 [500] /v2/lovecucu in 42.336µs for group v2
2021/04/27 19:20:28 [500] /v2/lovecucu in 67.767µs
*/

import (
	"gen"
	"log"
	"net/http"
	"time"
)

func onlyForV2() gen.HandlerFunc {
	return func(c *gen.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	engine := gen.Default()
	engine.GET("/", func(c *gen.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gen</h1>\n")
	})

	v2 := engine.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	{
		v2.GET("/hello/:name", func(c *gen.Context) {
			// expect /hello/lovecucu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	engine.Run(":8080")
}

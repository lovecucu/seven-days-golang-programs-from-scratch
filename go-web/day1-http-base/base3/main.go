package main

// $ curl http://localhost:8080/
// URL.Path = "/"
// curl http://localhost:8080/hello -H "a:b"
// 404 NOT FOUND: /hello
// $ curl -X POST http://localhost:8080/hello -H "a:b"
// Header["A"] = ["b"]
// Header["User-Agent"] = ["curl/7.64.1"]
// Header["Accept"] = ["*/*"]

import (
	"fmt"
	"gen"
	"log"
	"net/http"
)

func main() {
	engine := gen.New()
	engine.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	})
	engine.POST("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})
	log.Fatal(http.ListenAndServe(":8080", engine))
}

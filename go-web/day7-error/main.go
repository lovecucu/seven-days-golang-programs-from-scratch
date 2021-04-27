package main

/*
$ curl http://localhost:8080/
<html>
    <link rel="stylesheet" href="/assets/css/lovecucu.css">
    <p>geektutu.css is loaded</p>
</html>

$ curl http://localhost:8080/panic
{"message":"Internal Server Error"}

server的日志：
2021/04/27 20:50:03 runtime error: index out of range [100] with length 1
Traceback:
        /usr/local/Cellar/go/1.14.3/libexec/src/runtime/panic.go:969
        /usr/local/Cellar/go/1.14.3/libexec/src/runtime/panic.go:88
        /Users/bilibili/go/goproject/src/github.com/lovecucu/seven-days-golang-programs-from-scratch/go-web/day7-error/main.go:39
        /Users/bilibili/go/goproject/src/github.com/lovecucu/seven-days-golang-programs-from-scratch/go-web/day7-error/gen/context.go:43
        /Users/bilibili/go/goproject/src/github.com/lovecucu/seven-days-golang-programs-from-scratch/go-web/day7-error/gen/recovery.go:37
        /Users/bilibili/go/goproject/src/github.com/lovecucu/seven-days-golang-programs-from-scratch/go-web/day7-error/gen/context.go:43
        /Users/bilibili/go/goproject/src/github.com/lovecucu/seven-days-golang-programs-from-scratch/go-web/day7-error/gen/logger.go:15
        /Users/bilibili/go/goproject/src/github.com/lovecucu/seven-days-golang-programs-from-scratch/go-web/day7-error/gen/context.go:43
        /Users/bilibili/go/goproject/src/github.com/lovecucu/seven-days-golang-programs-from-scratch/go-web/day7-error/gen/router.go:102
        /Users/bilibili/go/goproject/src/github.com/lovecucu/seven-days-golang-programs-from-scratch/go-web/day7-error/gen/gen.go:124
        /usr/local/Cellar/go/1.14.3/libexec/src/net/http/server.go:2808
        /usr/local/Cellar/go/1.14.3/libexec/src/net/http/server.go:1896
        /usr/local/Cellar/go/1.14.3/libexec/src/runtime/asm_amd64.s:1374
*/

import (
	"fmt"
	"gen"
	"net/http"
	"time"
)

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%2d-%2d", year, month, day)
}

func main() {
	engine := gen.Default()

	engine.GET("/", func(c *gen.Context) {
		c.String(http.StatusOK, "Hello Geektutu\n")
	})
	// index out of range for testing Recovery()
	engine.GET("/panic", func(c *gen.Context) {
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[100]) // 索引越界会触发panic
	})

	engine.Run(":8080")
}

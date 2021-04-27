package main

/*
$ curl http://localhost:8080/
<html>
    <link rel="stylesheet" href="/assets/css/lovecucu.css">
    <p>geektutu.css is loaded</p>
</html>

$ curl http://localhost:8080/date

<html>
<body>
    <p>hello, gee</p>
    <p>Date: 2019- 8-17</p>
</body>
</html>


$ curl http://localhost:8080/students

<html>
<body>
    <p>hello, gee</p>

    <p>0: Geektutu is 20 years old</p>

    <p>1: Jack is 22 years old</p>

</body>
</html>
*/

import (
	"fmt"
	"gen"
	"html/template"
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
	engine.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	engine.LoadHTMLGlob("templates/*")
	engine.Static("/assets", "./static")

	stu1 := &student{Name: "Lovecucu", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	engine.GET("/", func(c *gen.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	engine.GET("/students", func(c *gen.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", gen.H{
			"title":  "gee",
			"stuArr": [2]*student{stu1, stu2},
		})
	})
	engine.GET("/date", func(c *gen.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", gen.H{
			"title": "gee",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	engine.Run(":8080")
}

package main

import (
	"fmt"
	"github.com/H-kang-better/msgo"
	"html/template"
	"log"
	"net/http"
)

type User struct {
	Name      string   `xml:"name" json:"name" msg:"required"`
	Age       int      `xml:"age" json:"age"`
	Addresses []string `xml:"addresses" json:"addresses"`
}

func main() {
	//http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
	//	fmt.Fprintln(writer, "hello mszlu.com")
	//})
	//err := http.ListenAndServe(":8111", nil)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//engine := msgo.New()
	//g := engine.Group("user")
	//g.Get("/hello", func(ctx *msgo.Context) {
	//	fmt.Fprintln(ctx.W, "hello mszlu.com")
	//})
	//g.Get("/hello/*", func(ctx *msgo.Context) {
	//	fmt.Fprintln(ctx.W, "hello* mszlu.com")
	//})
	//g.Post("/info", func(ctx *msgo.Context) {
	//	fmt.Fprintln(ctx.W, "info mszlu.com")
	//})
	//
	//order := engine.Group("order")
	//order.Get("/get/:id", func(ctx *msgo.Context) {
	//	fmt.Fprintln(ctx.W, "get order mszlu.com")
	//})
	//engine.Run()
	engine := msgo.New()
	g := engine.Group("user")
	g.Use(func(next msgo.HandlerFunc) msgo.HandlerFunc {
		return func(ctx *msgo.Context) {
			fmt.Println("pre handle")
			next(ctx)
			fmt.Println("post handle")
		}
	})
	g.Post("/hello/11", func(ctx *msgo.Context) {
		fmt.Fprintln(ctx.W, "post hello mszlu.com")
	})
	g.Get("/hello", func(ctx *msgo.Context) {
		fmt.Fprintln(ctx.W, "get hello mszlu.com")
	}, Log)
	// test 页面渲染功能的路由
	// test function HTML、HTMLTemplate
	g.Get("/HTML", func(ctx *msgo.Context) {
		fmt.Println("HTML")
		ctx.HTML(http.StatusOK, "<h1>你好 码神之路</h1>")
	})
	g.Get("/HTMLTemplate", func(ctx *msgo.Context) {
		fmt.Println("HTMLTemplate")
		ctx.HTMLTemplate("login.html", template.FuncMap{}, "", "tpl/index.html", "tpl/login.html", "tpl/header.html")
	})
	// test func HTMLTemplateGlob
	g.Get("/HTMLTemplateGlob", func(ctx *msgo.Context) {
		fmt.Println("HTMLTemplateGlob")
		ctx.HTMLTemplateGlob("login.html", template.FuncMap{}, "", "tpl/*.html")
	})
	// test func template 提前加载模板，比上面的加载方式简单
	engine.LoadTemplate("tpl/*.html") // 提前将模板加载到内存中
	g.Get("/template", func(ctx *msgo.Context) {
		err := ctx.Template("login.html", "")
		if err != nil {
			log.Println(err)
		}
	})
	// test func JSON
	g.Get("/json", func(ctx *msgo.Context) {
		_ = ctx.JSON(http.StatusOK, &User{
			Name: "码神之路",
		})
	})
	// test func XML
	g.Get("/xml", func(ctx *msgo.Context) {
		user := &User{
			Name: "码神之路",
		}
		_ = ctx.XML(http.StatusOK, user)
	})
	// test func File
	g.Get("/excel", func(ctx *msgo.Context) {
		ctx.File("tpl/test.xlsx")
	})
	// test func FileAttachment
	g.Get("/FileAttachment", func(ctx *msgo.Context) {
		ctx.FileAttachment("tpl/test.xlsx", "aaa")
	})
	// test func FileFromFS
	g.Get("/FileFromFS", func(ctx *msgo.Context) {
		//ctx.FileAttachment("tpl/test.xlsx", "哈哈.xlsx")
		ctx.FileFromFS("test.xlsx", http.Dir("tpl"))
	})
	// test func Redirect
	g.Get("/Redirect", func(ctx *msgo.Context) {
		ctx.Redirect(http.StatusFound, "/user/HTML")
	})
	// test func String
	g.Get("/String", func(ctx *msgo.Context) {
		ctx.String(http.StatusOK, "%s 是由 %s 制作 \n", "goweb框架", "码神之路")

	})
	// test func get value from url
	g.Get("/add", func(ctx *msgo.Context) {
		id := ctx.GetQuery("id")
		name, _ := ctx.GetQueryArray("id")
		country := ctx.GetDefaultQuery("country", "China")
		fmt.Println(id)
		fmt.Println(name)
		fmt.Println(country)
		ctx.String(http.StatusOK, "%s 是由 %s 制作 \n", "goweb框架", "码神之路")

	})
	// test func QueryMap
	g.Get("/QueryMap", func(ctx *msgo.Context) {
		userMap := ctx.QueryMap("user")
		fmt.Println(userMap)
		ctx.String(http.StatusOK, "%s 是由 %s 制作 \n", "goweb框架", "码神之路")
	})
	// test func PostForm
	g.Post("/PostForm", func(ctx *msgo.Context) {
		m, _ := ctx.GetPostForm("name")
		fmt.Println(m)
		ctx.JSON(http.StatusOK, m)
	})
	// test func FormFile
	g.Post("/FormFile", func(ctx *msgo.Context) {
		file, err := ctx.FormFile("file")
		if err != nil {
			log.Println(err)
		}
		err = ctx.SaveUploadedFile(file, "./upload/test.png")
		if err != nil {
			log.Println(err)
		}
	})
	// test func DealJson
	g.Post("/DealJson", func(ctx *msgo.Context) {
		ctx.DisallowUnknownFields = true
		ctx.IsValidate = true
		user := make([]User, 0)
		err := ctx.DealJson(&user)
		if err == nil {
			ctx.JSON(http.StatusOK, user)
		} else {
			log.Println(err)
		}
	})

	engine.Run()
}

func Log(next msgo.HandlerFunc) msgo.HandlerFunc {
	return func(ctx *msgo.Context) {
		fmt.Println("打印请求参数")
		next(ctx)
		fmt.Println("返回执行时间")
	}
}

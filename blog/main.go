package main

import (
	"fmt"
	"github.com/H-kang-better/msgo"
)

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
	g.Get("/hello/11", func(ctx *msgo.Context) {
		fmt.Fprintln(ctx.W, "get hello mszlu.com")
	}, Log)
	engine.Run()
}

func Log(next msgo.HandlerFunc) msgo.HandlerFunc {
	return func(ctx *msgo.Context) {
		fmt.Println("打印请求参数")
		next(ctx)
		fmt.Println("返回执行时间")
	}
}

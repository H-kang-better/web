package msgo

import (
	"fmt"
	"github.com/H-kang-better/msgo/render"
	"html/template"
	"log"
	"net/http"
	"strings"
	"sync"
)

const ANY = "ANY"

type HandlerFunc func(ctx *Context)

type MiddlewareFunc func(handlerFunc HandlerFunc) HandlerFunc

type routerGroup struct {
	groupName          string
	handlerFuncMap     map[string]map[string]HandlerFunc
	middlewaresFuncMap map[string]map[string][]MiddlewareFunc
	handlerMethodMap   map[string][]string //支持不同的路由请求
	treeNode           *treeNode
	middlewares        []MiddlewareFunc
}

type router struct {
	groups []*routerGroup
}

type Engine struct {
	*router
	funcMap    template.FuncMap
	HTMLRender render.HTMLRender
	pool       sync.Pool
}

func (r *routerGroup) Use(middlewares ...MiddlewareFunc) {
	r.middlewares = append(r.middlewares, middlewares...)
}

func (r *routerGroup) methodHandle(name string, method string, h HandlerFunc, ctx *Context) {
	//前置通用中间件
	if r.middlewares != nil {
		for _, middlewareFunc := range r.middlewares {
			h = middlewareFunc(h)
		}
	}
	//组路由级别
	middlewareFuncs := r.middlewaresFuncMap[name][method]
	if middlewareFuncs != nil {
		for _, middlewareFunc := range middlewareFuncs {
			h = middlewareFunc(h)
		}
	}
	h(ctx)
}

//func (r *routerGroup) Add(name string, handlerFunc HandlerFunc) {
//	r.handlerFuncMap[name] = handlerFunc
//}

func (r *routerGroup) handle(name string, method string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	_, ok := r.handlerFuncMap[name]
	if !ok {
		r.handlerFuncMap[name] = make(map[string]HandlerFunc)
		r.middlewaresFuncMap[name] = make(map[string][]MiddlewareFunc)
	}
	r.handlerFuncMap[name][method] = handlerFunc

	r.handlerMethodMap[method] = append(r.handlerMethodMap[method], name)
	r.middlewaresFuncMap[name][method] = append(r.middlewaresFuncMap[name][method], middlewareFunc...)
	r.treeNode.Put(name)
}

func (r *routerGroup) Any(name string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	r.handle(name, "ANY", handlerFunc, middlewareFunc...)
}

func (r *routerGroup) Get(name string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodGet, handlerFunc, middlewareFunc...)
}
func (r *routerGroup) Post(name string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodPost, handlerFunc, middlewareFunc...)
}
func (r *routerGroup) Delete(name string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodDelete, handlerFunc, middlewareFunc...)
}
func (r *routerGroup) Put(name string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodPut, handlerFunc, middlewareFunc...)
}
func (r *routerGroup) Patch(name string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodPatch, handlerFunc, middlewareFunc...)
}
func (r *routerGroup) Options(name string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodOptions, handlerFunc, middlewareFunc...)
}
func (r *routerGroup) Head(name string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodHead, handlerFunc, middlewareFunc...)
}

func (r *router) Group(name string) *routerGroup {
	g := &routerGroup{
		groupName:          name,
		handlerFuncMap:     make(map[string]map[string]HandlerFunc),
		middlewaresFuncMap: make(map[string]map[string][]MiddlewareFunc),
		handlerMethodMap:   make(map[string][]string),
		treeNode:           &treeNode{name: "/", children: make([]*treeNode, 0)},
	}
	r.groups = append(r.groups, g)
	return g
}

func New() *Engine {
	engine := &Engine{
		router:     &router{},
		HTMLRender: render.HTMLRender{},
	}
	engine.pool.New = func() any {
		return engine.allocateContext()
	}
	return engine
}

func (e *Engine) allocateContext() any {
	return &Context{engine: e}
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

// LoadTemplate 以 LoadTemplateGlob 方式加载所有模板
func (e *Engine) LoadTemplate(pattern string) {
	t := template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
	e.SetHtmlTemplate(t)
}

func (e *Engine) SetHtmlTemplate(t *template.Template) {
	e.HTMLRender = render.HTMLRender{Template: t}
}

func (e *Engine) httpRequestHandle(ctx *Context) {
	method := ctx.R.Method
	for _, group := range e.groups {
		routerName := SubStringLast(ctx.R.URL.Path, "/"+group.groupName)
		node := group.treeNode.Get(routerName)
		if node != nil && node.isEnd {
			// 路由匹配上了
			if handler, ok := group.handlerFuncMap[node.routerName][ANY]; ok {
				group.methodHandle(node.routerName, ANY, handler, ctx)
				return
			}

			if handler, ok := group.handlerFuncMap[node.routerName][ctx.R.Method]; ok {
				group.methodHandle(node.routerName, ctx.R.Method, handler, ctx)
				return
			}
			ctx.W.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(ctx.W, "%s %s not allowed \n", ctx.R.RequestURI, method)
			return
		}
	}
	ctx.W.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(ctx.W, "%s  not found \n", ctx.R.RequestURI)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := e.pool.Get().(*Context)
	ctx.W = w
	ctx.R = r
	e.httpRequestHandle(ctx)
	e.pool.Put(ctx)

}

func SubStringLast(str string, substr string) string {
	//先查找有没有
	index := strings.Index(str, substr)
	if index == -1 {
		return ""
	}
	return str[index+len(substr):]
}

func (e *Engine) Run() {
	//groups := e.router.groups
	//for _, g := range groups {
	//	for name, handle := range g.handlerFuncMap {
	//		http.HandleFunc("/"+g.groupName+name, handle)
	//	}
	//}
	http.Handle("/", e)
	err := http.ListenAndServe(":8111", nil)
	if err != nil {
		log.Fatal(err)
	}
}

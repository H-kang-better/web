package msgo

import (
	"github.com/H-kang-better/msgo/render"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

type Context struct {
	W      http.ResponseWriter
	R      *http.Request
	engine *Engine
}

func (c *Context) Render(code int, r render.Render) error {
	err := r.Render(c.W)
	c.W.WriteHeader(code)
	return err
}

// HTML 不支持模板的形式
func (c *Context) HTML(status int, html string) {
	c.Render(status, &render.HTML{IsTemplate: false, Data: html})
}

// HTMLTemplate 通过文件名称文件路径加载模板 ParseFiles(fileName...)
func (c *Context) HTMLTemplate(name string, funcMap template.FuncMap, data any, fileName ...string) {
	c.Render(http.StatusOK, &render.HTML{
		IsTemplate: true,
		Name:       name,
		Data:       data,
		Template:   c.engine.HTMLRender.Template,
	})
}

// HTMLTemplateGlob 通过 pattern 匹配，更简单
func (c *Context) HTMLTemplateGlob(name string, funcMap template.FuncMap, data any, pattern string) {
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	t := template.New(name)
	t.Funcs(funcMap)
	t, err := t.ParseGlob(pattern)
	if err != nil {
		log.Println(err)
		return
	}
	err = t.Execute(c.W, data)
	if err != nil {
		log.Println(err)
	}
}

// Template 启动的时候将所有模板加载到内存中，加快访问速度
func (c *Context) Template(name string, data any) error {
	return c.Render(http.StatusOK, &render.HTML{
		Name:       name,
		Template:   c.engine.HTMLRender.Template,
		Data:       data,
		IsTemplate: true})
}

// JSON 支持返回 json 格式
func (c *Context) JSON(status int, data any) error {
	return c.Render(status, &render.JSON{Data: data})
}

// XML 支持返回 xml 格式
func (c *Context) XML(status int, data any) error {
	return c.Render(status, &render.XML{Data: data})
}

// File 下载文件的需求，需要返回excel文件，word文件等等的
func (c *Context) File(filePath string) {
	http.ServeFile(c.W, c.R, filePath)
}

// FileAttachment 下载 filepath 路径下的文件后，将文件名修改为 filename
func (c *Context) FileAttachment(filepath, filename string) {
	if isASCII(filename) {
		c.W.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	} else {
		c.W.Header().Set("Content-Disposition", `attachment; filename*=UTF-8''`+url.QueryEscape(filename))
	}
	http.ServeFile(c.W, c.R, filepath)
}

// FileFromFS filepath 是相对于文件系统（fs）的相对路径
func (c *Context) FileFromFS(filepath string, fs http.FileSystem) {
	defer func(old string) {
		c.R.URL.Path = old
	}(c.R.URL.Path)

	c.R.URL.Path = filepath
	http.FileServer(fs).ServeHTTP(c.W, c.R)
}

// Redirect 重定向页面
func (c *Context) Redirect(status int, location string) error {
	return c.Render(status, &render.Redirect{
		Code:     status,
		Location: location,
		Request:  c.R,
	})
}

// String 用 value 填充 format 里面的值
func (c *Context) String(status int, format string, values ...any) error {
	return c.Render(status, &render.String{
		Format: format,
		Data:   values,
	})
}

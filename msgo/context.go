package msgo

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
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

// HTML 不支持模板的形式
func (c *Context) HTML(status int, html string) {
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.W.WriteHeader(status) // 设置返回状态
	_, err := c.W.Write([]byte(html))
	if err != nil {
		log.Println(err)
	}
}

// HTMLTemplate 通过文件名称文件路径加载模板 ParseFiles(fileName...)
func (c *Context) HTMLTemplate(name string, funcMap template.FuncMap, data any, fileName ...string) {
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	t := template.New(name)
	t.Funcs(funcMap)
	t, err := t.ParseFiles(fileName...)
	if err != nil {
		log.Println(err)
		return
	}

	err = t.Execute(c.W, data)
	if err != nil {
		log.Println(err)
	}
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
	c.W.Header().Set("Content-Type", "text/render; charset=utf-8")
	err := c.engine.HTMLRender.Template.ExecuteTemplate(c.W, name, data)
	return err
}

// JSON 支持返回 json 格式
func (c *Context) JSON(status int, data any) error {
	c.W.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.W.WriteHeader(status)
	rsp, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = c.W.Write(rsp)
	if err != nil {
		return err
	}
	return nil
}

// XML 支持返回 xml 格式
func (c *Context) XML(status int, data any) error {
	c.W.Header().Set("Content-Type", "application/xml; charset=utf-8")
	c.W.WriteHeader(status)
	//xmlData, err := xml.Marshal(data)
	//if err != nil {
	//	return err
	//}
	//_, err = c.W.Write(xmlData)
	//if err != nil {
	//	return err
	//}
	err := xml.NewEncoder(c.W).Encode(data)
	if err != nil {
		return err
	}
	return nil
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
func (c *Context) Redirect(status int, location string) {
	if (status < http.StatusMultipleChoices || status > http.StatusPermanentRedirect) && status != http.StatusCreated {
		panic(fmt.Sprintf("Cannot redirect with status code %d", status))
	}
	http.Redirect(c.W, c.R, location, status)
}

// String 用 value 填充 format 里面的值
func (c *Context) String(status int, format string, values ...any) (err error) {
	plainContentType := "text/plain; charset=utf-8"
	c.W.Header().Set("Content-Type", plainContentType)
	c.W.WriteHeader(status)
	if len(values) > 0 {
		_, err = fmt.Fprintf(c.W, format, values...)
		return
	}
	_, err = c.W.Write(StringToBytes(format))
	return
}

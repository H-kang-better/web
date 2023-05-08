package render

import (
	"html/template"
	"net/http"
)

type HTMLRender struct {
	Template *template.Template
}

type HTMLData any

type HTML struct {
	Template   *template.Template
	Name       string
	Data       HTMLData
	IsTemplate bool // 是否有模板
}

var htmlContentType = "text/html; charset=utf-8"

func (h *HTML) Render(w http.ResponseWriter) error {
	h.WriteContentType(w)
	if !h.IsTemplate {
		_, err := w.Write([]byte(h.Data.(string)))
		return err
	}
	err := h.Template.ExecuteTemplate(w, h.Name, h.Data)
	return err
}

func (h *HTML) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, htmlContentType)
}

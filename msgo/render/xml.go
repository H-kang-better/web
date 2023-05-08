package render

import (
	"encoding/xml"
	"net/http"
)

type XML struct {
	Data any
}

var xmlContentType = "application/xml; charset=utf-8"

func (x *XML) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, xmlContentType)
}

func (x *XML) Render(w http.ResponseWriter) error {
	x.WriteContentType(w)
	return xml.NewEncoder(w).Encode(x.Data)
}

package render

import (
	"fmt"
	"github.com/H-kang-better/msgo/internal/bytesconv"
	"net/http"
)

type String struct {
	Format string
	Data   []any
}

var plainContentType = "text/plain; charset=utf-8"

func (s *String) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, plainContentType)
}

func (s *String) Render(w http.ResponseWriter, code int) error {
	writeContentType(w, plainContentType)
	w.WriteHeader(code)
	if len(s.Data) > 0 {
		_, err := fmt.Fprintf(w, s.Format, s.Data...)
		return err
	}
	_, err := w.Write(bytesconv.StringToBytes(s.Format))
	return err
}

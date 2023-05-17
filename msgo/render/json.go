package render

import (
	"encoding/json"
	"net/http"
)

type JSON struct {
	Data any
}

var jsonContentType = "application/json; charset=utf-8"

func (j *JSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

func (j *JSON) Render(w http.ResponseWriter, code int) error {
	j.WriteContentType(w)
	w.WriteHeader(code)
	jsonBytes, err := json.Marshal(j.Data)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}

package render

import "net/http"

type Render interface {
	Bytes() ([]byte, error)
	Render(w http.ResponseWriter) error
	ContentType() string
	SetData(v interface{})
}

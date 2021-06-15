package render

import (
	"context"
	"net/http"
)

type Render interface {
	Bytes() ([]byte, error)
	Render(w http.ResponseWriter) error
	ContentType() string
	SetData(v interface{})
	SetHeader(k, v string)
	SetContext(ctx context.Context)
	SetError(err error)
}

package render

import (
	"fmt"
	"net/http"
)

type Text struct {
	content string
	Status  int
}

func (s Text) SetData(v interface{}) {
	s.content = fmt.Sprintf("%v", v)
}

func (s Text) Swagger() string {
	return "{}"
}

func (s Text) Bytes() ([]byte, error) {
	return []byte(s.content), nil
}

func (s Text) Render(w http.ResponseWriter) error {
	w.Header().Set("content-type", s.ContentType())
	if s.Status == 0 {
		s.Status = http.StatusOK
	}
	w.WriteHeader(s.Status)
	body, _ := s.Bytes()
	w.Write(body)
	return nil
}

func (s Text) ContentType() string {
	return "text/html"
}

package render

import (
	"context"
	"fmt"
	"net/http"
)

type Text struct {
	Content string
	Status  int
}

func (s *Text) SetHeader(k, v string) {
	panic("implement me")
}

func (s *Text) SetContext(ctx context.Context) {
	panic("implement me")
}

func (s *Text) SetError(err error) {
	panic("implement me")
}

func (s *Text) SetData(v interface{}) {
	s.Content = fmt.Sprintf("%v", v)
}

func (s *Text) Bytes() ([]byte, error) {
	return []byte(s.Content), nil
}

func (s *Text) Render(w http.ResponseWriter) error {
	w.Header().Set("content-type", s.ContentType())
	if s.Status == 0 {
		s.Status = http.StatusOK
	}
	w.WriteHeader(s.Status)
	body, _ := s.Bytes()
	w.Write(body)
	return nil
}

func (s *Text) ContentType() string {
	return "text/html"
}

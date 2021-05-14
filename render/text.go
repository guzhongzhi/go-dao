package render

import "net/http"

type Text struct {
	Content string
	Status  int
}

func (s Text) Bytes() ([]byte, error) {
	return []byte(s.Content), nil
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

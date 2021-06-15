package render

import (
	"context"
	"encoding/json"
	"net/http"
)

func NewJSON() Render {
	return &JSON{}
}

type JSON struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Status  int         `json:"status"`
}

func (s *JSON) SetHeader(k, v string) {
	panic("implement me")
}

func (s *JSON) SetContext(ctx context.Context) {
	panic("implement me")
}

func (s *JSON) SetError(err error) {
	panic("implement me")
}

func (s *JSON) SetData(v interface{}) {
	s.Data = v
}

func (s *JSON) Bytes() ([]byte, error) {
	b, _ := json.Marshal(map[string]interface{}{
		"data":    s.Data,
		"message": s.Message,
		"code":    s.Code,
		"status":  s.Status,
	})
	return b, nil
}

func (s *JSON) Render(w http.ResponseWriter) error {
	w.Header().Add("content-type", s.ContentType())
	w.WriteHeader(s.Status)
	body, err := s.Bytes()
	if err != nil {
		return err
	}
	w.Write(body)
	return nil
}

func (s *JSON) ContentType() string {
	return "application/json"
}

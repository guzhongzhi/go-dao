package render

import (
	"encoding/json"
	"github.com/guzhongzhi/gmicro/utils"
	"net/http"
	"reflect"
)

type JSON struct {
	data    interface{}
	Message string
	Code    int
	Status  int
}

func (s JSON) SetData(v interface{}) {
	s.data = v
}

func (s JSON) Swagger() string {
	data := map[string]interface{}{}
	t := reflect.TypeOf(data)
	js := utils.MapToJSON(t)
	body, _ := json.MarshalIndent(js, "", "    ")
	return string(body)
}

func (s JSON) Bytes() ([]byte, error) {
	b, _ := json.Marshal(map[string]interface{}{
		"data":    s.data,
		"message": s.Message,
		"code":    s.Code,
		"status":  s.Status,
	})
	return b, nil
}

func (s JSON) Render(w http.ResponseWriter) error {
	w.Header().Add("content-type", s.ContentType())
	w.WriteHeader(s.Status)
	body, err := s.Bytes()
	if err != nil {
		return err
	}
	w.Write(body)
	return nil
}

func (s JSON) ContentType() string {
	return "application/json"
}

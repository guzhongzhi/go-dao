package server

import (
	"fmt"
	"github.com/go-playground/form"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/guzhongzhi/gmicro/render"
	"net/http"
	"reflect"
)

type Router interface {
	HandlePath(mesh string, pathPattern string, call interface{}) Router
	Handle(mesh string, path runtime.Pattern, call interface{}) Router
	SetTagName(v string)
}

func NewRouter(mux *runtime.ServeMux) Router {
	s := &router{
		mux:     mux,
		decoder: form.NewDecoder(),
	}
	s.decoder.SetTagName("json")
	return s
}

type router struct {
	mux     *runtime.ServeMux
	decoder *form.Decoder
}

func (s *router) SetTagName(v string) {
	s.decoder.SetTagName(v)
}

func (s *router) HandlePath(mesh string, path string, call interface{}) Router {
	s.mux.HandlePath(mesh, path, s.handler(mesh, path, call))
	return s
}

func (s *router) Handle(mesh string, path runtime.Pattern, call interface{}) Router {
	s.mux.Handle(mesh, path, s.handler(mesh, path.String(), call))
	return s
}

func (s *router) buildCallParams(callType reflect.Type, r *http.Request, pathParams map[string]string) ([]reflect.Value, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	values := r.Form
	for k, v := range pathParams {
		values.Set(k, v)
	}

	if callType.NumIn() == 0 {
		return make([]reflect.Value, 0), nil
	}
	inType := callType.In(0)
	newIn := reflect.New(inType)

	in := newIn.Interface()
	err = s.decoder.Decode(in, values)
	if err != nil {
		return nil, err
	}
	if inType.Kind() == reflect.Ptr {
		in = newIn.Interface()
	} else {
		in = newIn.Elem().Interface()
	}
	return []reflect.Value{reflect.ValueOf(in)}, nil
}

func (s *router) handler(mesh string, path string, call interface{}) runtime.HandlerFunc {
	c := reflect.TypeOf(call)
	if c.Kind() != reflect.Func {
		panic(fmt.Sprintf("the %s:%s call of http handle must bu a func", mesh, path))
	}

	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		params, err := s.buildCallParams(c, r, pathParams)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("form parse error: " + err.Error()))
			return
		}

		v := reflect.ValueOf(call)
		rsp := v.Call(params)
		if len(rsp) == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("there is no response of the call '%s'", v.String())))
			return
		}
		rr, ok := rsp[0].Interface().(render.Render)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("the response of the call of '%s' must be a render", v.String())))
			return
		}
		rr.Render(w)
	}
}

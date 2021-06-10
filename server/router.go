package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/form"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/guzhongzhi/gmicro/render"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

type Router interface {
	HandlePath(mesh string, pathPattern string, call interface{}) Router
	Handle(mesh string, path runtime.Pattern, call interface{}) Router
	SetTagName(v string)
	Swagger() *Swagger
	EnableSwaggerJSON(handlerFunc ...runtime.HandlerFunc)
}

func NewRouter(mux *runtime.ServeMux) Router {
	s := &router{
		mux:     mux,
		decoder: form.NewDecoder(),
		paths:   make([]Path, 0),
		swagger: NewSwagger(),
	}
	s.SetTagName("json")
}

type Path struct {
	Method   string
	Path     string
	Params   interface{}
	Response interface{}
}

type router struct {
	tagName string
	mux     *runtime.ServeMux
	decoder *form.Decoder
	paths   []Path
	swagger *Swagger
}

func (s *router) Swagger() *Swagger {
	return s.swagger
}

func (s *router) SetTagName(v string) {
	s.tagName = v
	s.decoder.SetTagName(v)
	s.swagger.tagName = v
}

func (s *router) EnableSwaggerJSON(handlerFunc ...runtime.HandlerFunc) {
	s.mux.HandlePath("OPTIONS", "/", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT")
		w.Header().Set("Access-Control-Allow-Origin", "*")
	})
	s.mux.HandlePath("GET", "/swagger.json", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		body := s.swagger.ToJSON()
		w.Write([]byte(body))
	})
}

func (s *router) HandlePath(mesh string, path string, call interface{}) Router {
	s.mux.HandlePath(mesh, path, s.handler(mesh, path, call))
	return s
}

func (s *router) Handle(mesh string, path runtime.Pattern, call interface{}) Router {
	s.mux.Handle(mesh, path, s.handler(mesh, path.String(), call))
	return s
}

func (s *router) buildCallParams(ctx context.Context, callType reflect.Type, r *http.Request, pathParams map[string]string) ([]reflect.Value, error) {
	if callType.NumIn() == 0 {
		return make([]reflect.Value, 0), nil
	}

	if callType.NumIn() == 1 {
		return []reflect.Value{reflect.ValueOf(ctx)}, nil
	}

	inType := callType.In(1)
	newIn := reflect.New(inType)
	in := newIn.Interface()
	var err error

	if strings.Index(r.Header.Get("content-type"), "application/json") != -1 {
		var body []byte
		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(body, in)
	} else {
		err := r.ParseForm()
		if err != nil {
			return nil, err
		}
		values := r.Form
		for k, v := range pathParams {
			values.Set(k, v)
		}
		err = s.decoder.Decode(in, values)
	}
	if err != nil {
		return nil, err
	}

	if inType.Kind() == reflect.Ptr {
		in = newIn.Elem().Interface()
	}
	return []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(in)}, nil
}

func (s *router) handler(mesh string, path string, call interface{}, opts ...RouterHandlerOption) runtime.HandlerFunc {
	c := reflect.TypeOf(call)

	if c.Kind() != reflect.Func {
		panic(fmt.Sprintf("the %s:%s call of http handle must bu a func", mesh, path))
	}
	routerOptions := NewRouterHandlerOptions(opts ...)
	if routerOptions.ResponseWrapper != nil {
		_, ok := routerOptions.ResponseWrapper.(render.Render)
		if !ok {
			panic(fmt.Sprintf("invalid response wrapper for '%s', method='%s'", path, mesh))
		}
	}
	s.swagger.AddHandler(mesh, path, c.In(1), c.Out(0), routerOptions)

	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		params, err := s.buildCallParams(r.Context(), c, r, pathParams)

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
		if len(rsp) >= 2 {
			data := rsp[1].Interface()
			rr.SetData(data)
		}
		rr.Render(w)
	}
}

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/form"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/guzhongzhi/gmicro/render"
	"github.com/guzhongzhi/gmicro/utils"
	"net/http"
	"reflect"
)

type Context struct {
	context.Context
}

func (s Context) Request() *http.Request {
	v := s.Value("r_context")
	if v == nil {
		return nil
	}
	return v.(*http.Request)
}

type Router interface {
	HandlePath(mesh string, pathPattern string, call interface{}) Router
	Handle(mesh string, path runtime.Pattern, call interface{}) Router
	SetTagName(v string)
}

func NewRouter(mux *runtime.ServeMux) Router {
	s := &router{
		mux:     mux,
		decoder: form.NewDecoder(),
		paths:   make([]Path, 0),
	}
	s.SetTagName("json")
	mux.HandlePath("GET", "/routers", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		h := `<style>
.method-title {
	background:#000000;
	padding:10px 4px;
	font-weight: bold;
	font-size: 16px;
	color: #FFF;
}
.method {
	padding:10px;
}
.method-params-title, .method-out-title {
	padding:4px 0px;

}
.method-params-body, .method-out-body {
    margin-top:10px;
	background:#e1e1e1;
	padding:10px;
}
</style>`
		for _, p := range s.paths {
			h += "<div class='method'>"
			h += fmt.Sprintf("<div class='method-title'> %s  %s<br></div>", p.Method, p.Path)
			h += fmt.Sprintf("<div class='method-params'><div class='method-params-title'>Request:</div><pre class='method-params-body'>%s</pre></div>", p.Params)
			h += fmt.Sprintf("<div class='method-out'><div class='method-out-title'>Response:</div><pre class='method-out-body'>%s</pre></div>", p.Response)
			h += "</div>"
		}
		w.Header().Set("content-type", "text/html")
		w.Write([]byte(h))
	})
	return s
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
}

func (s *router) SetTagName(v string) {
	s.tagName = v
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
	ctx := Context{
		r.Context(),
	}

	if callType.NumIn() == 1 {
		return []reflect.Value{reflect.ValueOf(ctx)}, nil
	}

	inType := callType.In(1)
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
	return []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(in)}, nil
}

func (s *router) createInJSON(callType reflect.Type) interface{} {
	if callType.NumIn() <= 1 {
		return "{}"
	}
	inType := callType.In(1)
	fields := s.loopType(inType)
	js, _ := json.MarshalIndent(fields, "", "    ")
	return js
}

func (s *router) createOutJSON(callType reflect.Type) interface{} {
	if callType.NumOut() <= 0 {
		return "{}"
	}
	inType := callType.Out(1)
	fields := s.loopType(inType)
	js, _ := json.MarshalIndent(fields, "", "    ")
	return js
}

func (s *router) loopType(inType reflect.Type) interface{} {
	return utils.LoopType(s.tagName, inType)
}

func (s *router) handler(mesh string, path string, call interface{}) runtime.HandlerFunc {
	c := reflect.TypeOf(call)
	if c.Kind() != reflect.Func {
		panic(fmt.Sprintf("the %s:%s call of http handle must bu a func", mesh, path))
	}
	s.paths = append(s.paths, Path{
		Method:   mesh,
		Path:     path,
		Params:   s.createInJSON(c),
		Response: s.createOutJSON(c),
	})

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
		if len(rsp) >= 2 {
			data := rsp[1].Interface()
			rr.SetData(data)
		}
		rr.Render(w)
	}
}

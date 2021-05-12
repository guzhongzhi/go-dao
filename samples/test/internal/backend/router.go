package backend

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/guzhongzhi/gmicro/server"
	"google.golang.org/grpc"
	"net/http"
)

func NewRegister() server.Registry {
	return &Registry{
	}
}

type Registry struct {
}

func (s *Registry) Register(mux *runtime.ServeMux, server *grpc.Server) {
	mux.HandlePath("GET", "/", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.Write([]byte("eee"))
	})
}

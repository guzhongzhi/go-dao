package application

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/guzhongzhi/gmicro/server"
	"github.com/guzhongzhi/gmicro/test/api"
	"google.golang.org/grpc"
	"net/http"
)

func NewRegister(SubEffect api.SubEffectServiceServer) server.Registry {
	return &Registry{
		SubEffect: SubEffect,
	}
}

type Registry struct {
	SubEffect api.SubEffectServiceServer
}

func (s *Registry) Register(mux *runtime.ServeMux, server *grpc.Server) {
	if server != nil {
		api.RegisterSubEffectServiceServer(server, s.SubEffect)
	}
	api.RegisterSubEffectServiceHandlerServer(context.Background(), mux, s.SubEffect)
	mux.HandlePath("GET", "/a", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.Write([]byte("eee"))
	})
}

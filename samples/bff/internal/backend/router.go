package backend

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/guzhongzhi/gmicro/server"
	"google.golang.org/grpc"
)

func NewRegister() server.Registry {
	return &Registry{
		user: &User{},
	}
}

type Registry struct {
	user *User
}

func (s *Registry) Register(mux *runtime.ServeMux, server *grpc.Server, router server.Router) {
	router.SetTagName("json")
	router.
		HandlePath("POST", "/user", s.user.Create).
		HandlePath("GET", "/user/{id}", s.user.Get).
		HandlePath("PUT", "/user/{id}", s.user.Update).
		HandlePath("DELETE", "/user/{id}", s.user.Delete).
		HandlePath("PATCH", "/user/{id}", s.user.Update)
}

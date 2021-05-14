package backend

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/guzhongzhi/gmicro/client"
	"github.com/guzhongzhi/gmicro/logger"
	"github.com/guzhongzhi/gmicro/server"
	"github.com/guzhongzhi/gmicro/test/api"
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
	mux.HandlePath("GET", "/b", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		c, err := client.NewGRPCClient("test", "127.0.0.1", 9000, nil)
		fmt.Println(err)
		err = c.Callback(func(conn *grpc.ClientConn, log logger.SuperLogger) error {
			c := api.NewSubEffectServiceClient(conn)
			in := &api.UpsertRequest{}
			_, err := c.Create(context.Background(), in)
			return err
		})
		fmt.Println(err)
		in := &api.UpsertRequest{}
		rsp := &api.UpsertResponse{}
		err = c.Call(context.Background(), "/api.SubEffectService/Create", in, rsp)
		fmt.Println(err)
		w.Write([]byte("eee"))
	})
}

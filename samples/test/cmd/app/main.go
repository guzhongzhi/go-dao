package main

import (
	"context"
	"fmt"
	"github.com/guzhongzhi/gmicro/config"
	"github.com/guzhongzhi/gmicro/console"
	"github.com/guzhongzhi/gmicro/logger"
	"github.com/guzhongzhi/gmicro/server"
	"github.com/guzhongzhi/gmicro/test/internal/infrastructure"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net/http"
	"os"
	"path"
)

func main() {

	basePath := path.Dir(path.Dir(os.Args[0]))
	app := console.NewApp("salad-effect", "1.0", basePath)
	app.Action = func(ctx *cli.Context) error {
		env := ctx.String("env")
		cfgPath := ctx.String("config")

		cfg := &infrastructure.Bootstrap{}
		err := config.LoadConfigByFiles(cfgPath, env, cfg, logger.Default())
		if err != nil {
			panic(err)
			os.Exit(1)
		}

		serverConfig := server.NewConfig(server.GRPCAddrOption(cfg.ServerConfig().GRPC.Addr),
			server.GRPCServerOption(grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
				md, _ := metadata.FromIncomingContext(ctx)
				fmt.Println("md: ", md)
				return handler(ctx, req)
			})),
			server.HTTPPluginsOption(cfg.ServerConfig().HTTP.Plugins),
			server.HTTPAddrOption(cfg.ServerConfig().HTTP.Addr),
			server.HTTPHandlerOption(func(h http.Handler, logger logger.SuperLogger) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Println("wrapper in main")
					h.ServeHTTP(w, r)
				})
			}),
		)

		server, _, err := initApp(cfg, logger.Default(), serverConfig)
		if err != nil {
			panic(err)
		}
		return server.Serve()
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

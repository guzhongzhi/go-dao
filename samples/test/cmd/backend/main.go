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
)

func main() {
	os.Setenv("SALAD_SERVER.GRPC.DISABLED", "true")
	os.Setenv("SALAD_SERVER.HTTP.ADDR", "0.0.0.0:8009")
	envPrefix := "SALAD"
	cfg := infrastructure.NewBootstrap()

	consoleCfg := console.NewConfig("salad-effect", "1.0", "", cfg)
	consoleCfg.EnvPrefix = envPrefix
	app := console.New(consoleCfg)
	app.App().Action = func(ctx *cli.Context) error {
		env := ctx.String("env")
		cfgPath := ctx.String("config")

		err := config.LoadConfigFiles(cfgPath, env, cfg, logger.Default(), envPrefix)
		if err != nil {
			panic(err)
			os.Exit(1)
		}

		serverConfig := server.ApplyOptions(cfg.ServerConfig(), server.GRPCAddrOption(cfg.ServerConfig().GRPC.Addr),
			server.GRPCServerOption(grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
				md, _ := metadata.FromIncomingContext(ctx)
				fmt.Println("md: ", md)
				return handler(ctx, req)
			})),
			server.HTTPPluginsOption(cfg.ServerConfig().HTTP.Plugins),
			server.HTTPAddrOption(cfg.ServerConfig().HTTP.Addr),
			server.HTTPHandlerOption(func(h http.Handler, logger logger.SuperLogger) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Println("wrapper backend in main")
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

	err := app.App().Run(os.Args)
	if err != nil {
		panic(err)
	}
}

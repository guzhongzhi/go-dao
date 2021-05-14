//+build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/google/wire"
	"github.com/guzhongzhi/gmicro/logger"
	"github.com/guzhongzhi/gmicro/server"
	"github.com/guzhongzhi/gmicro/samples/bff/internal/frontend"
	"github.com/guzhongzhi/gmicro/samples/bff/internal/infrastructure"
)

//initApp init kratos application.
//go:generate kratos t wire
func initApp(
	cfg *infrastructure.Bootstrap,
	l logger.SuperLogger,
	serverOptions *server.Config,
) (*server.Server, func(), error) {
	panic(wire.Build(
		server.NewServer,
		frontend.ProviderSet,
	))
}

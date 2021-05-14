// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"github.com/guzhongzhi/gmicro/logger"
	"github.com/guzhongzhi/gmicro/server"
	"github.com/guzhongzhi/gmicro/test/internal/backend"
	"github.com/guzhongzhi/gmicro/test/internal/infrastructure"
)

// Injectors from wire.go:

//initApp init kratos application.
//go:generate kratos t wire
func initApp(cfg *infrastructure.Bootstrap, l logger.SuperLogger, serverOptions *server.Config) (*server.Server, func(), error) {
	registry := backend.NewRegister()
	serverServer := server.NewServer(serverOptions, registry, l)
	return serverServer, func() {
	}, nil
}
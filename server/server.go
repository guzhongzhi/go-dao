package server

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pinguo-icc/salad-effect/internal/infrastructure/logger"
	"github.com/pinguo-icc/salad-effect/internal/infrastructure/server/middleware"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Registry interface {
	Register(mux *runtime.ServeMux, server *grpc.Server)
}

type Server struct {
	config     *Config
	register   Registry
	grpcServer *grpc.Server
	httpServer *http.Server
	logger     logger.SuperLogger
	sysSig     chan os.Signal
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Infof("start to stop servers")
	s.grpcServer.Stop()
	s.httpServer.Close()
	return nil
}

func (s *Server) syscall() error {
	signal.Notify(s.sysSig, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-s.sysSig
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			s.Stop(context.Background())
			time.Sleep(time.Second * 2)
			s.logger.Infof("servers are stopped")
			os.Exit(0)
		case syscall.SIGHUP:
		default:
			return nil
		}
	}
	return nil
}

func (s *Server) Serve() error {
	grpcListener, err := net.Listen("tcp", s.config.GRPC.Addr)
	if err != nil {
		panic(err)
	}
	s.grpcServer = grpc.NewServer(s.config.GRPC.Options...)

	mux := runtime.NewServeMux()
	s.register.Register(mux, s.grpcServer)

	var httpHandler = http.Handler(mux)
	for _, name := range s.config.HTTP.Plugins {
		if fn, ok := middleware.Middlewares[name]; ok {
			httpHandler = fn(httpHandler, s.logger)
		}
	}

	s.httpServer = &http.Server{
		Addr:    s.config.HTTP.Addr,
		Handler: httpHandler,
	}

	go func() {
		err := s.grpcServer.Serve(grpcListener)
		if err != nil {
			s.logger.Infof("grpc error: %s", err.Error())
		}
		s.logger.Infof("grpc server stopped")

	}()
	go func() {
		err := s.httpServer.ListenAndServe()
		if err != nil {
			s.logger.Info(err.Error())
		}
	}()

	s.logger.Infof("grpc listen: %s", s.config.GRPC.Addr)
	s.logger.Infof("http listen: %s", s.config.HTTP.Addr)

	return s.syscall()
}

func NewServer(config *Config, register Registry, logger logger.SuperLogger) *Server {
	return &Server{
		config:   config,
		register: register,
		logger:   logger,
		sysSig:   make(chan os.Signal, 1),
	}
}

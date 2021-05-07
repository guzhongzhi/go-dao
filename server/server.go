package server

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"time"
)

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

type Option func(opts *Config)

type Config struct {
	GRPCAddr          string
	GRPCServerOptions []grpc.ServerOption
	Timeout           time.Duration
	HTTPAddr          string
	Handlers          []http.Handler
}

func DefaultOptions() *Config {
	return &Config{
		GRPCAddr:          "0.0.0.0:9000",
		GRPCServerOptions: make([]grpc.ServerOption, 0),
		HTTPAddr:          "0.0.0.0:8000",
		Handlers:          make([]http.Handler, 0),
	}
}

func NewOptions(opts ...Option) *Config {
	o := DefaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func GRPCAddrOption(v string) Option {
	return func(opts *Config) {
		opts.GRPCAddr = v
	}
}

func HTTPAddrOption(v string) Option {
	return func(opts *Config) {
		opts.HTTPAddr = v
	}
}

type Registry interface {
	Register(mux *runtime.ServeMux, server *grpc.Server)
}

type Server struct {
	config     *Config
	register   Registry
	grpcServer *grpc.Server
	httpServer *http.Server
}

func (s *Server) Stop(ctx context.Context) error {
	s.grpcServer.Stop()
	return s.httpServer.Close()
}

func (s *Server) Serve() error {
	grpcListener, err := net.Listen("tcp", s.config.GRPCAddr)
	if err != nil {
		panic(err)
	}
	s.grpcServer = grpc.NewServer(s.config.GRPCServerOptions...)

	mux := runtime.NewServeMux()
	s.register.Register(mux, s.grpcServer)

	s.httpServer = &http.Server{
		Addr:    s.config.HTTPAddr,
		Handler: mux,
	}

	sig := make(chan error, 1)
	go func() {
		err := s.grpcServer.Serve(grpcListener)
		sig <- err
	}()
	go func() {
		err := s.httpServer.ListenAndServe()
		sig <- err
	}()

	fmt.Println(fmt.Sprintf("grpc listen: %s", s.config.GRPCAddr))
	fmt.Println(fmt.Sprintf("http listen: %s", s.config.HTTPAddr))
	select {
	case err = <-sig:
		panic(err)
	default:

	}
	return nil
}

func NewServer(options *Config, register Registry) *Server {
	return &Server{
		config:   options,
		register: register,
	}
}

package server

import (
	"github.com/guzhongzhi/gmicro/server/middleware"
	"google.golang.org/grpc"
)

type Option func(opts *Config)

type httpConfig struct {
	Addr     string
	Plugins  []string
	Handlers []middleware.Middleware
}

type grpcConfig struct {
	Addr    string
	Options []grpc.ServerOption
}

type Config struct {
	HTTP *httpConfig
	GRPC *grpcConfig
}

func DefaultConfig() *Config {
	return &Config{
		GRPC: &grpcConfig{
			Addr:    "0.0.0.0:9000",
			Options: make([]grpc.ServerOption, 0),
		},
		HTTP: &httpConfig{
			Addr:     "0.0.0.0:8000",
			Handlers: make([]middleware.Middleware, 0),
			Plugins:  make([]string, 0),
		},
	}
}

func NewConfig(opts ...Option) *Config {
	o := DefaultConfig()
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func GRPCServerOption(opt ...grpc.ServerOption) Option {
	return func(opts *Config) {
		opts.GRPC.Options = append(opts.GRPC.Options, opt...)
	}
}

func GRPCAddrOption(v string) Option {
	return func(opts *Config) {
		opts.GRPC.Addr = v
	}
}

func HTTPAddrOption(v string) Option {
	return func(opts *Config) {
		opts.HTTP.Addr = v
	}
}

func HTTPPluginsOption(v []string) Option {
	return func(opts *Config) {
		opts.HTTP.Plugins = v
	}
}

func HTTPHandlerOption(v middleware.Middleware) Option {
	return func(opts *Config) {
		opts.HTTP.Handlers = append(opts.HTTP.Handlers, v)
	}
}

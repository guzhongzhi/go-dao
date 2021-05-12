package server

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/guzhongzhi/gmicro/logger"
	"github.com/guzhongzhi/gmicro/server/middleware"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
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

	for _, fn := range s.config.HTTP.Handlers {
		httpHandler = fn(httpHandler, s.logger)
	}

	s.httpServer = &http.Server{
		Addr:    s.config.HTTP.Addr,
		Handler: httpHandler,
	}

	go func() {
		ip, port := parseAddr(grpcListener.Addr().String())
		s.config.GRPC.Addr = fmt.Sprintf("%s:%s", ip, port)
		s.logger.Infof("grpc listen: %s", s.config.GRPC.Addr)

		err := s.grpcServer.Serve(grpcListener)
		if err != nil {
			s.logger.Infof("grpc error: %s", err.Error())
		}
		s.logger.Infof("grpc server stopped")

	}()
	go func() {
		httpListener, err := net.Listen("tcp", s.config.HTTP.Addr)
		if err != nil {
			panic(err)
		}
		ip, port := parseAddr(httpListener.Addr().String())
		s.config.HTTP.Addr = fmt.Sprintf("%s:%s", ip, port)
		s.logger.Infof("http listen: %s", s.config.HTTP.Addr)

		err = s.httpServer.Serve(httpListener)
		if err != nil {
			s.logger.Info(err.Error())
		}
	}()

	return s.syscall()
}

func parseAddr(addr string) (string, string) {
	temp := strings.Split(addr, ":")
	port := temp[len(temp)-1]
	s := strings.Join(temp[:len(temp)-1], ":")
	if strings.Index(s, ":") != -1 {
		ips := localIP()
		if len(ips) > 0 {
			s = ips[0]
		}
	}
	return s, port
}

func localIP() []string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	regex, err := regexp.Compile("^(10|172|192.168)")
	if err != nil {
		panic(err)
	}

	ips := make([]string, 0)

	for _, address := range addrs {
		ipnet, ok := address.(*net.IPNet);
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}
		if ipnet.IP.To4() != nil && regex.MatchString(ipnet.IP.String()) {
			ips = append(ips, ipnet.IP.String())
		}
	}
	return ips
}

func NewServer(config *Config, register Registry, logger logger.SuperLogger) *Server {
	return &Server{
		config:   config,
		register: register,
		logger:   logger,
		sysSig:   make(chan os.Signal, 1),
	}
}

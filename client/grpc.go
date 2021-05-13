package client

import (
	"context"
	"fmt"
	"github.com/guzhongzhi/gmicro/logger"
	"google.golang.org/grpc"
)

type GRCPClient interface {
	Callback(call func(conn *grpc.ClientConn, logger logger.SuperLogger) error) error
	Call(ctx context.Context, method string, in interface{}, out interface{}, opts ...grpc.CallOption) error
}

func NewGRPCClient(serviceName string, port int, l logger.SuperLogger) (*grpcClient, error) {
	if l == nil {
		l = logger.Default()
	}
	c := &grpcClient{
		serviceName: serviceName,
		port:        port,
		logger:      l,
	}
	return c, nil
}

type grpcClient struct {
	serviceName string
	port        int
	logger      logger.SuperLogger
}

func (s *grpcClient) connect() (*grpc.ClientConn, error) {
	addr := s.serviceName
	if s.port > 0 {
		addr = fmt.Sprintf("%s:%v", addr, s.port)
	}
	return grpc.Dial(addr, grpc.WithInsecure())
}

func (s *grpcClient) Callback(call func(conn *grpc.ClientConn, logger logger.SuperLogger) error) error {
	conn, err := s.connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	return call(conn, s.logger)
}

func (s *grpcClient) Call(ctx context.Context, method string, in interface{}, out interface{}, opts ...grpc.CallOption) error {
	conn, err := s.connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	s.logger.Debugf("start to send grpc call to '%s', method='%s'", s.serviceName, method)
	err = conn.Invoke(ctx, method, in, out, opts...)
	s.logger.Debugf("end send grpc call to '%s', method='%s'", s.serviceName, method)
	if err != nil {
		return err
	}
	return nil
}

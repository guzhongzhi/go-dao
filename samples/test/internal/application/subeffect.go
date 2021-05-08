package application

import "github.com/guzhongzhi/gmicro/test/api"

func NewSubEffectServer() api.SubEffectServiceServer {
	return &SubEffectServer{}
}

type SubEffectServer struct {
	api.UnimplementedSubEffectServiceServer
}

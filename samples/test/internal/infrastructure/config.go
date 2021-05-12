package infrastructure

import (
	"github.com/guzhongzhi/gmicro/server"
)

func NewBootstrap() *Bootstrap {
	return &Bootstrap{
		Server: server.DefaultConfig(),
	}
}

// Bootstrap 所有配置实例
type Bootstrap struct {
	Server *server.Config `mapstructure:"server"`
}

func (s *Bootstrap) ServerConfig() *server.Config {
	if s.Server == nil {
		return server.DefaultConfig()
	}
	return s.Server
}

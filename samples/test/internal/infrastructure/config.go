package infrastructure

import (
	"github.com/guzhongzhi/gmicro/server"
)

// Bootstrap 所有配置实例
type Bootstrap struct {
	Server *server.Config `mapstructure:"server"`
}

func (s *Bootstrap) ServerConfig() *server.Config {
	return s.Server
}

package application

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewSubEffectServer,
	NewRegister,
)

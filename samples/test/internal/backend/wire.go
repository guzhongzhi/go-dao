package backend

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewRegister,
)

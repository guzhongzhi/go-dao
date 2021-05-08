package middleware

import (
	logger2 "github.com/pinguo-icc/salad-effect/internal/infrastructure/logger"
	"net/http"
)

type Middleware func(h http.Handler, logger logger2.SuperLogger) http.Handler

var Middlewares = make(map[string]Middleware)

func init() {
	Middlewares["cors"] = AllowCORS
	Middlewares["logger"] = Logger
}

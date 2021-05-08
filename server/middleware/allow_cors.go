package middleware

import (
	logger2 "github.com/pinguo-icc/salad-effect/internal/infrastructure/logger"
	"net/http"
)

func AllowCORS(h http.Handler, superLogger logger2.SuperLogger) http.Handler {
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

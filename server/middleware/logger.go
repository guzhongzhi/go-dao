package middleware

import (
	logger2 "github.com/guzhongzhi/gmicro/logger"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var logger = logrus.New()

func init() {
	logger.SetFormatter(&logrus.JSONFormatter{})
}

func Logger(h http.Handler, superLogger logger2.SuperLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		go superLogger.Infof("%v %s %s %s %s",
			time.Now().Format("2006-01-02 15:04:05.999999999 -0700 MST"),
			r.Method, r.URL.Path,
			r.URL.Query().Encode(), r.UserAgent())
		h.ServeHTTP(w, r)
	})
}

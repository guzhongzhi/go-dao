package logger

import "github.com/sirupsen/logrus"

func Default() SuperLogger {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.SetReportCaller(false)
	log.SetFormatter(&logrus.JSONFormatter{})
	return log
}

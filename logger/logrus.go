package logger

import "github.com/sirupsen/logrus"

func Default() SuperLogger {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	return log
}

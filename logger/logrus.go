package logger

import "github.com/sirupsen/logrus"

func Default() SuperLogger {
	return logrus.New()
}

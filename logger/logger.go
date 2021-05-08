package logger

type Loggerf interface {
	Debugf(format string, args ...interface{})

	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})

	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

type Logger interface {
	Debug(args ...interface{})

	Info(args ...interface{})
	Warn(args ...interface{})

	Error(args ...interface{})
	Fatal(args ...interface{})
}

type CloserLogger interface {
	Close() error
}

type SuperLogger interface {
	Logger
	Loggerf
}

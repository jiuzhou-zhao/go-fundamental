package interfaces

type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

type EmptyLogger struct {
}

func (logger *EmptyLogger) Debug(v ...interface{}) {

}

func (logger *EmptyLogger) Debugf(format string, v ...interface{}) {

}
func (logger *EmptyLogger) Info(v ...interface{}) {

}
func (logger *EmptyLogger) Infof(format string, v ...interface{}) {

}
func (logger *EmptyLogger) Warn(v ...interface{}) {

}
func (logger *EmptyLogger) Warnf(format string, v ...interface{}) {

}
func (logger *EmptyLogger) Error(v ...interface{}) {

}
func (logger *EmptyLogger) Errorf(format string, v ...interface{}) {

}
func (logger *EmptyLogger) Fatal(v ...interface{}) {

}
func (logger *EmptyLogger) Fatalf(format string, v ...interface{}) {

}

package logger

type NoopLogger struct{}

func (n NoopLogger) Info(msg string, keyVals ...interface{})     {}
func (n NoopLogger) Debug(msg string, keyVals ...interface{})    {}
func (n NoopLogger) Warn(msg string, keyVals ...interface{})     {}
func (n NoopLogger) Error(msg string, keyVals ...interface{})    {}
func (n NoopLogger) Fatal(msg string, keyVals ...interface{})    {}
func (n NoopLogger) Errorf(msg string, keyVals ...interface{})   {}
func (n NoopLogger) Warningf(msg string, keyVals ...interface{}) {}
func (n NoopLogger) Infof(msg string, keyVals ...interface{})    {}
func (n NoopLogger) Debugf(msg string, keyVals ...interface{})   {}
func (n NoopLogger) With(keyVals ...interface{}) Logger          { return n }

func NewNopLogger() Logger {
	return NoopLogger{}
}

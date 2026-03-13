// Package zap is a stub of go.uber.org/zap for use in analysistest testdata.
package zap

// Field is a zap field.
type Field struct{}

// Logger is a stub of zap.Logger.
type Logger struct{}

// NewNop returns a no-op Logger.
func NewNop() *Logger { return &Logger{} }

func (l *Logger) Debug(msg string, fields ...Field)  {}
func (l *Logger) Info(msg string, fields ...Field)   {}
func (l *Logger) Warn(msg string, fields ...Field)   {}
func (l *Logger) Error(msg string, fields ...Field)  {}
func (l *Logger) DPanic(msg string, fields ...Field) {}
func (l *Logger) Panic(msg string, fields ...Field)  {}
func (l *Logger) Fatal(msg string, fields ...Field)  {}

// Sugar returns a SugaredLogger.
func (l *Logger) Sugar() *SugaredLogger { return &SugaredLogger{} }

// String creates a string Field.
func String(key, val string) Field { return Field{} }

// SugaredLogger is a stub of zap.SugaredLogger.
type SugaredLogger struct{}

func (s *SugaredLogger) Debug(args ...interface{})              {}
func (s *SugaredLogger) Info(args ...interface{})               {}
func (s *SugaredLogger) Warn(args ...interface{})               {}
func (s *SugaredLogger) Error(args ...interface{})              {}
func (s *SugaredLogger) DPanic(args ...interface{})             {}
func (s *SugaredLogger) Panic(args ...interface{})              {}
func (s *SugaredLogger) Fatal(args ...interface{})              {}
func (s *SugaredLogger) Debugf(tmpl string, args ...interface{}) {}
func (s *SugaredLogger) Infof(tmpl string, args ...interface{})  {}
func (s *SugaredLogger) Warnf(tmpl string, args ...interface{})  {}
func (s *SugaredLogger) Errorf(tmpl string, args ...interface{}) {}
func (s *SugaredLogger) DPanicf(tmpl string, args ...interface{}) {}
func (s *SugaredLogger) Panicf(tmpl string, args ...interface{})  {}
func (s *SugaredLogger) Fatalf(tmpl string, args ...interface{})  {}

package logger

//go:generate stringer -type=Level
type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
	Critical
)

func (l *Logger) Debug(message string, fields ...Fieldser) {
	l.Log(Debug, message, fields...)
}

func (l *Logger) Info(message string, fields ...Fieldser) {
	l.Log(Info, message, fields...)
}

func (l *Logger) Warn(message string, fields ...Fieldser) {
	l.Log(Warn, message, fields...)
}

func (l *Logger) Error(message string, fields ...Fieldser) {
	l.Log(Error, message, fields...)
}

func (l *Logger) Critical(message string, fields ...Fieldser) {
	l.Log(Critical, message, fields...)
}

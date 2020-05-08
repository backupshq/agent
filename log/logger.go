package log

const Debug = 4
const Info = 3
const Warn = 2
const Error = 1

type Logger struct {
	level  int
	writer Writer
}

type Writer interface {
	Write(level int, message string)
}

func (l *Logger) Debug(message string) {
	if l.level >= Debug {
		l.writer.Write(Debug, message)
	}
}
func (l *Logger) Info(message string) {
	if l.level >= Info {
		l.writer.Write(Info, message)
	}
}
func (l *Logger) Warn(message string) {
	if l.level >= Warn {
		l.writer.Write(Warn, message)
	}
}
func (l *Logger) Error(message string) {
	if l.level >= Error {
		l.writer.Write(Error, message)
	}
}

package log

const Debug = 5
const Info = 4
const Warn = 3
const Error = 2
const Fatal = 1

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
func (l *Logger) Fatal(message string) {
	if l.level >= Fatal {
		l.writer.Write(Fatal, message)
	}
}

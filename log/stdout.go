package log

import (
	"fmt"
)

type StdoutWriter struct {
}

func (w *StdoutWriter) Write(level int, message string) {
	fmt.Println(withTime(withLevel(level, message)))
}

func CreateStdoutLogger(level int) *Logger {
	return &Logger{
		level:  level,
		writer: &StdoutWriter{},
	}
}

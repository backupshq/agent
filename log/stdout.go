package log

import (
	"fmt"
	"strings"
)

type StdoutWriter struct {
}

func (w *StdoutWriter) Write(level int, message string) {
	lines := strings.Split(message, "\n")
	for _, line := range lines {
		fmt.Println(withTime(withLevel(level, line)))
	}
}

func CreateStdoutLogger(level int) *Logger {
	return &Logger{
		level:  level,
		writer: &StdoutWriter{},
	}
}

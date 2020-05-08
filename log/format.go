package log

import (
	"fmt"
	"time"
)

var levels = map[int]string{
	Debug: "DEBUG",
	Info:  "INFO",
	Warn:  "WARN",
	Error: "ERROR",
}

func levelLabel(level int) string {
	if level, ok := levels[level]; ok {
		return level
	}
	return ""
}

func withTime(message string) string {
	return fmt.Sprintf("%s %s", time.Now().UTC().Format(time.RFC3339), message)
}

func withLevel(level int, message string) string {
	label := fmt.Sprintf("[%s]", levelLabel(level))
	label += " "

	return fmt.Sprintf("%s %s", label[:7], message)
}

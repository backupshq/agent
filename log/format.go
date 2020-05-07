package log

import (
	"fmt"
	"time"
)

func levelLabel(level int) string {
	switch level {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	default:
		return ""
	}
}

func withTime(message string) string {
	return fmt.Sprintf("%s %s", time.Now().UTC().Format(time.RFC3339), message)
}

func withLevel(level int, message string) string {
	label := fmt.Sprintf("[%s]", levelLabel(level))
	label += " "

	return fmt.Sprintf("%s %s", label[:7], message)
}

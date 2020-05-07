package log

import (
	"testing"
)

type MockLog struct {
	level   int
	message string
}
type MockWriter struct {
	messages []MockLog
}

func (w *MockWriter) Write(level int, message string) {
	w.messages = append(w.messages, MockLog{level, message})
}

func testLog(l *Logger) {
	l.Debug("Test debug")
	l.Info("Test info")
	l.Warn("Test warn")
	l.Error("Test error")
}

func assertLog(t *testing.T, log MockLog, expectedMessage string) {
	if log.message != expectedMessage {
		t.Errorf("wrong message logged, expected %q but got %q", expectedMessage, log.message)
	}
}

func TestLogger(t *testing.T) {
	t.Run("Test debug", func(t *testing.T) {
		w := &MockWriter{}
		l := &Logger{
			level:  Debug,
			writer: w,
		}
		testLog(l)

		if len(w.messages) != 4 {
			t.Errorf("expected 4 messages, got %d", len(w.messages))
			return
		}

		assertLog(t, w.messages[0], "Test debug")
		assertLog(t, w.messages[1], "Test info")
		assertLog(t, w.messages[2], "Test warn")
		assertLog(t, w.messages[3], "Test error")
	})

	t.Run("Test info", func(t *testing.T) {
		w := &MockWriter{}
		l := &Logger{
			level:  Info,
			writer: w,
		}

		testLog(l)

		if len(w.messages) != 3 {
			t.Errorf("expected 3 messages, got %d", len(w.messages))
			return
		}

		assertLog(t, w.messages[0], "Test info")
		assertLog(t, w.messages[1], "Test warn")
		assertLog(t, w.messages[2], "Test error")
	})

	t.Run("Test warn", func(t *testing.T) {
		w := &MockWriter{}
		l := &Logger{
			level:  Warn,
			writer: w,
		}

		testLog(l)

		if len(w.messages) != 2 {
			t.Errorf("expected 2 messages, got %d", len(w.messages))
			return
		}

		assertLog(t, w.messages[0], "Test warn")
		assertLog(t, w.messages[1], "Test error")
	})

	t.Run("Test error", func(t *testing.T) {
		w := &MockWriter{}
		l := &Logger{
			level:  Error,
			writer: w,
		}

		testLog(l)

		if len(w.messages) != 1 {
			t.Errorf("expected 1 message, got %d", len(w.messages))
			return
		}

		assertLog(t, w.messages[0], "Test error")
	})
}

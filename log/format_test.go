package log

import (
	"testing"
)

func assertEquals(t *testing.T, expected string, actual string) {
	if actual != expected {
		t.Errorf("Expected %q but got %q", expected, actual)
	}
}

func TestWithLevel(t *testing.T) {
	t.Run("white space is padded", func(t *testing.T) {
		assertEquals(t, "[DEBUG] Hello world", withLevel(Debug, "Hello world"))
		assertEquals(t, "[INFO]  Hello world", withLevel(Info, "Hello world"))
		assertEquals(t, "[WARN]  Hello world", withLevel(Warn, "Hello world"))
		assertEquals(t, "[ERROR] Hello world", withLevel(Error, "Hello world"))
	})
}

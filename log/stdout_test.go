package log

import (
	"testing"
)

func TestStdoutLogger(t *testing.T) {
	t.Run("Logger has correct level", func(t *testing.T) {
		l := CreateStdoutLogger(Debug)
		if l.level != Debug {
			t.Errorf("Expected level to be %d but got %d", Debug, l.level)
		}

		l = CreateStdoutLogger(Error)
		if l.level != Error {
			t.Errorf("Expected level to be %d but got %d", Error, l.level)
		}
	})
}

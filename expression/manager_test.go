package expression

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestEvaluateSimple(t *testing.T) {
	cases := map[string]string{
		"foo":               "foo",
		"%'test'%":          "test",
		"one %'two'% three": "one two three",
	}

	for input, expected := range cases {
		t.Run("Evaluate", func(t *testing.T) {
			manager := CreateExpressionManager()

			output, err := manager.Evaluate(input, Context{})
			if err != nil {
				t.Errorf("got unexpected error %q", err)
				return
			}
			if !cmp.Equal(expected, output) {
				t.Errorf(cmp.Diff(expected, output))
			}
		})
	}
}

func TestEvaluateWithVariables(t *testing.T) {
	cases := map[string]string{
		"%one%":                            "1",
		"I have %one% dog and %two% cats.": "I have 1 dog and 2 cats.",
		"This is 90%% finished.":           "This is 90% finished.",
	}

	for input, expected := range cases {
		t.Run("Evaluate", func(t *testing.T) {
			manager := CreateExpressionManager()

			output, err := manager.Evaluate(input, Context{
				map[string]string{
					"one": "1",
					"two": "2",
				},
			})
			if err != nil {
				t.Errorf("got unexpected error %q", err)
				return
			}
			if !cmp.Equal(expected, output) {
				t.Errorf(cmp.Diff(expected, output))
			}
		})
	}
}

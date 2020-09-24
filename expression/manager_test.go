package expression

import (
	"github.com/google/go-cmp/cmp"
	"math"
	"strconv"
	"strings"
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
				map[string]func(args ...string) string{},
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

func TestEvaluateWithFunctions(t *testing.T) {
	cases := map[string]string{
		"%f()%":                     "default",
		"%f('')%":                   "default",
		"%f('test')%":               "grfg",
		"%f('test')%/%f('foo')%":    "grfg/sbb",
		"%f(f('test'))%":            "test",
		"%g(one)%":                  "1",
		"%g(two)%":                  "4",
		"%h(one, two, f('three'))%": "12guerr",
		"%h(h,h,h)%":                "hhh",
	}

	for input, expected := range cases {
		t.Run("Evaluate", func(t *testing.T) {
			manager := CreateExpressionManager()

			output, err := manager.Evaluate(input, Context{
				map[string]string{
					"one": "1",
					"two": "2",
					"h":   "h",
				},
				map[string]func(args ...string) string{
					"f": func(args ...string) string {
						if len(args) > 0 && len(args[0]) > 0 {
							return strings.Map(
								// rot 13
								func(r rune) rune {
									if r >= 'a' && r <= 'z' {
										if r >= 'm' {
											return r - 13
										} else {
											return r + 13
										}
									} else if r >= 'A' && r <= 'Z' {
										if r >= 'M' {
											return r - 13
										} else {
											return r + 13
										}
									}
									// Do nothing.
									return r
								},
								args[0],
							)
						}
						return "default"
					},
					"g": func(args ...string) string {
						if len(args) < 1 {
							return ""
						}
						num, err := strconv.Atoi(args[0])
						if err != nil {
							return err.Error()
						}
						return strconv.Itoa(int(math.Pow(float64(num), float64(num))))
					},
					"h": func(args ...string) string {
						output := ""

						for _, s := range args {
							output += s
						}

						return output
					},
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

func TestBadSyntax(t *testing.T) {
	cases := []string{
		"%£^&£$^&$^£(f('test', 'test'))%",
		"f('test')%",
		"%f('test')",
		"%f(%)%",
		"%‮f('test')%",
		"%%%%%",
		"%f('test'))%",
		"%f(('test')%",
		"%h('test',, 'test')%",
		"%h('test', ''test')%",
		"%h('test', \"test')%",
		"%h('test', test')%",
		"%test'%",
		"%one'%",
	}

	for _, input := range cases {
		t.Run("Parse", func(t *testing.T) {
			manager := CreateExpressionManager()

			_, err := manager.Parse(input)
			if err == nil {
				t.Errorf("expected error, got nil")
				return
			}
		})
	}
}

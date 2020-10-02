package expression

import (
	"errors"
	"fmt"
)

type Context struct {
	Variables map[string]string
	Functions map[string]func(args ...string) string
}

func (c *Context) getVariable(name string) (string, error) {
	if value, ok := c.Variables[name]; ok {
		return value, nil
	}

	return "", errors.New(fmt.Sprintf("Unknown variable '%s'", name))
}

func (c *Context) getFunction(name string) (func(args ...string) string, error) {
	if function, ok := c.Functions[name]; ok {
		return function, nil
	}

	defaultFunc := func(args ...string) string {
		return ""
	}

	return defaultFunc, errors.New(fmt.Sprintf("Unknown function '%s'", name))
}

package expression

import (
	"errors"
	"fmt"
)

type Context struct {
	variables map[string]string
}

func (c *Context) getVariable(name string) (string, error) {
	if value, ok := c.variables[name]; ok {
		return value, nil
	}

	return "", errors.New(fmt.Sprintf("Unknown variable '%s'", name))
}

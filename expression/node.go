package expression

import (
	"errors"
	"fmt"
)

type Node interface {
	Evaluate(context Context) (string, error)
}

type ParentNode struct {
	children []Node
}

func (n *ParentNode) Evaluate(context Context) (string, error) {
	output := ""

	for _, child := range n.children {
		out, err := child.Evaluate(context)
		if err != nil {
			return output, err
		}

		output += out
	}

	return output, nil
}

type ConstantNode struct {
	value string
}

func (n *ConstantNode) Evaluate(context Context) (string, error) {
	return n.value, nil
}

type VariableNode struct {
	name string
}

func (n *VariableNode) Evaluate(context Context) (string, error) {
	return context.getVariable(n.name)
}

type FunctionNode struct {
	name string
	args []Node
}

func (n *FunctionNode) Evaluate(context Context) (string, error) {
	function, err := context.getFunction(n.name)
	if err != nil {
		return "", err
	}

	var resolvedArgs []string
	for i, arg := range n.args {
		resolvedArg, err := arg.Evaluate(context)
		if err != nil {
			return "", errors.New(fmt.Sprintf("Error at argument %d: %s", i, err.Error()))
		}

		resolvedArgs = append(resolvedArgs, resolvedArg)
	}

	return function(resolvedArgs...), nil
}

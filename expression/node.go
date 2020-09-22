package expression

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

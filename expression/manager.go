package expression

type ExpressionManager struct {
	lexer  Lexer
	parser Parser
}

func (e *ExpressionManager) Parse(input string) (Node, error) {
	tokens, err := e.lexer.Tokenize(input)

	if err != nil {
		return nil, err
	}

	return e.parser.Parse(tokens)
}

func (e *ExpressionManager) Evaluate(input string, context Context) (string, error) {
	node, err := e.Parse(input)

	if err != nil {
		return "", err
	}

	return node.Evaluate(context)
}

func CreateExpressionManager() *ExpressionManager {
	return &ExpressionManager{
		Lexer{},
		Parser{},
	}
}

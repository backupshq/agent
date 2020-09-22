package expression

type ExpressionManager struct {
	lexer  Lexer
	parser Parser
}

func (e *ExpressionManager) Evaluate(input string, context Context) (string, error) {
	tokens, err := e.lexer.Tokenize(input)

	if err != nil {
		return "", err
	}

	node, err := e.parser.Parse(tokens)

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

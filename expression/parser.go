package expression

import (
	"errors"
	"fmt"
)

type Parser struct{}

func (p *Parser) Parse(tokens *TokenList) (Node, error) {
	var children []Node
	var err error

	for tokens.current != nil {
		if tokens.current.is(TypeLiteral, "") {
			children = append(children, &ConstantNode{tokens.current.Value})
			tokens.next()
			continue
		}

		err = tokens.expect(TypeExpressionBoundary, "%")
		if err != nil {
			return nil, err
		}
		parsedExpression, err := p.parseExpression(tokens)
		if err != nil {
			return nil, err
		}
		children = append(children, parsedExpression)
		err = tokens.expect(TypeExpressionBoundary, "%")
		if err != nil {
			return nil, err
		}
	}

	return &ParentNode{children}, nil
}

func (p *Parser) parseExpression(tokens *TokenList) (Node, error) {
	if tokens.current == nil {
		return nil, errors.New("Unexpected end of input.")
	}

	if tokens.current.is(TypeString, "") {
		value := tokens.current.Value
		tokens.next()

		return &ConstantNode{value}, nil
	}

	if tokens.current.is(TypeIdentifier, "") {
		identifierName := tokens.current.Value
		tokens.next()

		if tokens.current.is(TypeBracket, "(") {
			tokens.next()
			args, err := p.parseArguments(tokens)
			if err != nil {
				return nil, err
			}
			tokens.expect(TypeBracket, ")")

			return &FunctionNode{identifierName, args}, nil
		}

		return &VariableNode{identifierName}, nil
	}

	return nil, errors.New(fmt.Sprintf("At position %d: unexpected token '%s'", tokens.current.Position, tokens.current.Value))
}

func (p *Parser) parseArguments(tokens *TokenList) ([]Node, error) {
	var args []Node

	for !tokens.current.is(TypeBracket, ")") {
		parsedExpression, err := p.parseExpression(tokens)
		if err != nil {
			return nil, err
		}
		args = append(args, parsedExpression)
		if !tokens.current.is(TypeBracket, "") {
			tokens.expect(TypeComma, "")
		}
	}

	return args, nil
}

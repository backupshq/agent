package expression

import (
	"errors"
	"fmt"
)

// Anything outside of %
const TypeLiteral = "literal"

// %
const TypeExpressionBoundary = "expression_boundary"

// function or variable
const TypeIdentifier = "identifier"

// ( or )
const TypeBracket = "bracket"

// ,
const TypeComma = "comma"

// 'string'
const TypeString = "string"

type Token struct {
	TokenType string
	Position  int
	Value     string
}

func (t Token) is(tokenType, value string) bool {
	if value != "" && value != t.Value {
		return false
	}

	return tokenType == t.TokenType
}

type TokenList struct {
	tokens  []Token
	current *Token
	cursor  int
}

func CreateTokenList(tokens []Token) *TokenList {
	list := &TokenList{
		tokens,
		nil,
		-1,
	}
	list.next()

	return list
}

func (l *TokenList) next() {
	l.cursor++
	if l.cursor < len(l.tokens) {
		l.current = &l.tokens[l.cursor]
	} else {
		l.current = nil
	}
}

func (l *TokenList) expect(tokenType, value string) error {
	if l.current == nil {
		return errors.New(fmt.Sprintf("Unexpected end of input, expected '%s'.", tokenType))
	}

	if !l.current.is(tokenType, value) {
		msg := fmt.Sprintf("Unexpected token '%s', expected '%s'.", l.current.TokenType, tokenType)

		return errors.New(msg)
	}

	l.next()
	return nil
}

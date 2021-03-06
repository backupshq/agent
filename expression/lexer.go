package expression

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Lexer struct{}

func (l *Lexer) Tokenize(input string) (*TokenList, error) {
	cursor := 0
	var tokens []Token
	end := len([]byte(input))
	inExpression := false
	var err error
	var match []int

	regexLiteral := regexp.MustCompile(`^[^%]*(%%)?[^%]*`)
	regexIdentifier := regexp.MustCompile(`^[a-zA-Z][_a-zA-Z0-9]*`)
	regexString := regexp.MustCompile(`^'[^']*'`)

	for cursor < end {
		if !inExpression {
			match = regexLiteral.FindStringIndex(input[cursor:])
			if match != nil && match[1] != 0 {
				tokens = append(tokens, Token{TypeLiteral, cursor, strings.ReplaceAll(input[cursor:][0:match[1]], "%%", "%")})
				cursor = cursor + match[1]
				continue
			}
			tokens = append(tokens, Token{TypeExpressionBoundary, cursor, "%"})
			cursor++
			inExpression = true
			continue
		}

		if input[cursor] == '%' {
			tokens = append(tokens, Token{TypeExpressionBoundary, cursor, "%"})
			cursor++
			inExpression = false
			continue
		}

		if input[cursor] == '(' || input[cursor] == ')' {
			tokens = append(tokens, Token{TypeBracket, cursor, input[cursor : cursor+1]})
			cursor++
			continue
		}

		if input[cursor] == ',' {
			tokens = append(tokens, Token{TypeComma, cursor, ","})
			cursor++
			continue
		}

		// ignore whitespace
		if input[cursor] == ' ' {
			cursor++
			continue
		}

		match = regexIdentifier.FindStringIndex(input[cursor:])
		if match != nil {
			tokens = append(tokens, Token{TypeIdentifier, cursor, input[cursor:][match[0]:match[1]]})
			cursor = cursor + match[1]
			continue
		}

		match = regexString.FindStringIndex(input[cursor:])
		if match != nil {
			var matchedString string
			if match[1]-match[0] < 2 {
				matchedString = ""
			} else {
				matchedString = input[cursor:][match[0]+1 : match[1]-1]
			}
			tokens = append(tokens, Token{TypeString, cursor, matchedString})
			cursor = cursor + match[1]
			continue
		}

		currentRune, _ := utf8.DecodeRuneInString(input[cursor:])
		err = errors.New(fmt.Sprintf("At position %d: unexpected input %s", cursor, string(currentRune)))
		break
	}

	return CreateTokenList(tokens), err
}

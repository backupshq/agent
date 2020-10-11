package expression

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestTokenize(t *testing.T) {
	cases := map[string][]Token{
		"test": {
			{TypeLiteral, 0, "test"},
		},

		"test %'foo'%": {
			{TypeLiteral, 0, "test "},
			{TypeExpressionBoundary, 5, "%"},
			{TypeString, 6, "foo"},
			{TypeExpressionBoundary, 11, "%"},
		},

		"test %   'foo'%": {
			{TypeLiteral, 0, "test "},
			{TypeExpressionBoundary, 5, "%"},
			{TypeString, 9, "foo"},
			{TypeExpressionBoundary, 14, "%"},
		},

		"test %bar%": {
			{TypeLiteral, 0, "test "},
			{TypeExpressionBoundary, 5, "%"},
			{TypeIdentifier, 6, "bar"},
			{TypeExpressionBoundary, 9, "%"},
		},

		"test %run()%": {
			{TypeLiteral, 0, "test "},
			{TypeExpressionBoundary, 5, "%"},
			{TypeIdentifier, 6, "run"},
			{TypeBracket, 9, "("},
			{TypeBracket, 10, ")"},
			{TypeExpressionBoundary, 11, "%"},
		},

		"test %run('foo')%": {
			{TypeLiteral, 0, "test "},
			{TypeExpressionBoundary, 5, "%"},
			{TypeIdentifier, 6, "run"},
			{TypeBracket, 9, "("},
			{TypeString, 10, "foo"},
			{TypeBracket, 15, ")"},
			{TypeExpressionBoundary, 16, "%"},
		},

		"test %run( bar )%": {
			{TypeLiteral, 0, "test "},
			{TypeExpressionBoundary, 5, "%"},
			{TypeIdentifier, 6, "run"},
			{TypeBracket, 9, "("},
			{TypeIdentifier, 11, "bar"},
			{TypeBracket, 15, ")"},
			{TypeExpressionBoundary, 16, "%"},
		},

		"((([]as398&£JKDM": {
			{TypeLiteral, 0, "((([]as398&£JKDM"},
		},

		"%%": {
			{TypeLiteral, 0, "%"},
		},

		"100%%": {
			{TypeLiteral, 0, "100%"},
		},

		"100%% %'foo'% %% ": {
			{TypeLiteral, 0, "100% "},
			{TypeExpressionBoundary, 6, "%"},
			{TypeString, 7, "foo"},
			{TypeExpressionBoundary, 12, "%"},
			{TypeLiteral, 13, " % "},
		},

		"%f(a, b, 'c')%": {
			{TypeExpressionBoundary, 0, "%"},
			{TypeIdentifier, 1, "f"},
			{TypeBracket, 2, "("},
			{TypeIdentifier, 3, "a"},
			{TypeComma, 4, ","},
			{TypeIdentifier, 6, "b"},
			{TypeComma, 7, ","},
			{TypeString, 9, "c"},
			{TypeBracket, 12, ")"},
			{TypeExpressionBoundary, 13, "%"},
		},
	}

	for input, expectedTokens := range cases {
		t.Run("Tokenize", func(t *testing.T) {
			lexer := Lexer{}

			list, err := lexer.Tokenize(input)
			if err != nil {
				t.Errorf("got unexpected error %q", err)
				return
			}
			if !cmp.Equal(expectedTokens, list.tokens) {
				t.Errorf(cmp.Diff(expectedTokens, list.tokens))
			}
		})
	}
}

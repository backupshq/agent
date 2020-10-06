package expression

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestTokenize(t *testing.T) {
	cases := map[string][]Token{
		"test": []Token{
			Token{TypeLiteral, 0, "test"},
		},

		"test %'foo'%": []Token{
			Token{TypeLiteral, 0, "test "},
			Token{TypeExpressionBoundary, 5, "%"},
			Token{TypeString, 6, "foo"},
			Token{TypeExpressionBoundary, 11, "%"},
		},

		"test %   'foo'%": []Token{
			Token{TypeLiteral, 0, "test "},
			Token{TypeExpressionBoundary, 5, "%"},
			Token{TypeString, 9, "foo"},
			Token{TypeExpressionBoundary, 14, "%"},
		},

		"test %bar%": []Token{
			Token{TypeLiteral, 0, "test "},
			Token{TypeExpressionBoundary, 5, "%"},
			Token{TypeIdentifier, 6, "bar"},
			Token{TypeExpressionBoundary, 9, "%"},
		},

		"test %run()%": []Token{
			Token{TypeLiteral, 0, "test "},
			Token{TypeExpressionBoundary, 5, "%"},
			Token{TypeIdentifier, 6, "run"},
			Token{TypeBracket, 9, "("},
			Token{TypeBracket, 10, ")"},
			Token{TypeExpressionBoundary, 11, "%"},
		},

		"test %run('foo')%": []Token{
			Token{TypeLiteral, 0, "test "},
			Token{TypeExpressionBoundary, 5, "%"},
			Token{TypeIdentifier, 6, "run"},
			Token{TypeBracket, 9, "("},
			Token{TypeString, 10, "foo"},
			Token{TypeBracket, 15, ")"},
			Token{TypeExpressionBoundary, 16, "%"},
		},

		"test %run( bar )%": []Token{
			Token{TypeLiteral, 0, "test "},
			Token{TypeExpressionBoundary, 5, "%"},
			Token{TypeIdentifier, 6, "run"},
			Token{TypeBracket, 9, "("},
			Token{TypeIdentifier, 11, "bar"},
			Token{TypeBracket, 15, ")"},
			Token{TypeExpressionBoundary, 16, "%"},
		},

		"((([]as398&£JKDM": []Token{
			Token{TypeLiteral, 0, "((([]as398&£JKDM"},
		},

		"%%": []Token{
			Token{TypeLiteral, 0, "%"},
		},

		"100%%": []Token{
			Token{TypeLiteral, 0, "100%"},
		},

		"100%% %'foo'% %% ": []Token{
			Token{TypeLiteral, 0, "100% "},
			Token{TypeExpressionBoundary, 6, "%"},
			Token{TypeString, 7, "foo"},
			Token{TypeExpressionBoundary, 12, "%"},
			Token{TypeLiteral, 13, " % "},
		},

		"%f(a, b, 'c')%": []Token{
			Token{TypeExpressionBoundary, 0, "%"},
			Token{TypeIdentifier, 1, "f"},
			Token{TypeBracket, 2, "("},
			Token{TypeIdentifier, 3, "a"},
			Token{TypeComma, 4, ","},
			Token{TypeIdentifier, 6, "b"},
			Token{TypeComma, 7, ","},
			Token{TypeString, 9, "c"},
			Token{TypeBracket, 12, ")"},
			Token{TypeExpressionBoundary, 13, "%"},
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

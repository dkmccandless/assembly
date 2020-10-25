package lexer

import (
	"testing"

	"github.com/dkmccandless/assembly/token"
)

func TestNext(t *testing.T) {
	tests := []struct {
		input  string
		tokens []token.Token
	}{
		{
			input: `WHEREAS RESOLVED`,
			tokens: []token.Token{
				{token.WHEREAS, "WHEREAS"},
				{token.RESOLVED, "RESOLVED"},
				{token.EOF, ""},
			},
		},
		{
			input: `A Resolution Concerning Commentary

WHEREAS a resolution consisting entirely of comments has no effect: now, therefore, 

BE IT RESOLVED that this assembly takes no action.`,
			tokens: []token.Token{
				{token.WHEREAS, "WHEREAS"},
				{token.RESOLVED, "RESOLVED"},
				{token.EOF, ""},
			},
		},
	}
	for _, test := range tests {
		l := New(test.input)
		for _, want := range test.tokens {
			// Disregard comment tokens in input
			got := l.Next()
			for got.Typ == token.COMMENT {
				got = l.Next()
			}

			if got != want {
				t.Errorf("Next(%v): got %v, want %v", test.input, got, want)
			}
		}
	}
}

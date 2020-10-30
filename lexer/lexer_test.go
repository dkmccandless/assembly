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
			input: `()-`,
			tokens: []token.Token{
				{token.LPAREN, "("},
				{token.RPAREN, ")"},
				{token.DASH, "-"},
				{token.EOF, ""},
			},
		},
		{
			input: `WHEREAS RESOLVED`,
			tokens: []token.Token{
				{token.WHEREAS, "WHEREAS"},
				{token.RESOLVED, "RESOLVED"},
				{token.EOF, ""},
			},
		},
		{
			input: `-1 2 3000000000000`,
			tokens: []token.Token{
				{token.NUMERAL, "-1"},
				{token.NUMERAL, "2"},
				{token.NUMERAL, "3000000000000"},
				{token.EOF, ""},
			},
		},
		{
			input: "negative zero one two three four five six seven eight nine ten eleven twelve thirteen fourteen fifteen sixteen seventeen eighteen nineteen twenty thirty forty fifty sixty seventy eighty ninety hundred thousand million billion trillion quadrillion quintillion",
			tokens: []token.Token{
				{token.NEGATIVE, "negative"},
				{token.ZERO, "zero"},
				{token.ONES, "one"},
				{token.ONES, "two"},
				{token.ONES, "three"},
				{token.ONES, "four"},
				{token.ONES, "five"},
				{token.ONES, "six"},
				{token.ONES, "seven"},
				{token.ONES, "eight"},
				{token.ONES, "nine"},
				{token.VIGESIMAL, "ten"},
				{token.VIGESIMAL, "eleven"},
				{token.VIGESIMAL, "twelve"},
				{token.VIGESIMAL, "thirteen"},
				{token.VIGESIMAL, "fourteen"},
				{token.VIGESIMAL, "fifteen"},
				{token.VIGESIMAL, "sixteen"},
				{token.VIGESIMAL, "seventeen"},
				{token.VIGESIMAL, "eighteen"},
				{token.VIGESIMAL, "nineteen"},
				{token.TENS, "twenty"},
				{token.TENS, "thirty"},
				{token.TENS, "forty"},
				{token.TENS, "fifty"},
				{token.TENS, "sixty"},
				{token.TENS, "seventy"},
				{token.TENS, "eighty"},
				{token.TENS, "ninety"},
				{token.HUNDRED, "hundred"},
				{token.POWER, "thousand"},
				{token.POWER, "million"},
				{token.POWER, "billion"},
				{token.POWER, "trillion"},
				{token.POWER, "quadrillion"},
				{token.POWER, "quintillion"},
				{token.EOF, ""},
			},
		},
		{
			input: "negative three (-3)",
			tokens: []token.Token{
				{token.NEGATIVE, "negative"},
				{token.ONES, "three"},
				{token.LPAREN, "("},
				{token.NUMERAL, "-3"},
				{token.RPAREN, ")"},
				{token.EOF, ""},
			},
		},
		{
			input: `""`,
			tokens: []token.Token{
				{token.STRING, ""},
				{token.EOF, ""},
			},
		},
		{
			input: `"Greetings, Assembly."`,
			tokens: []token.Token{
				{token.STRING, "Greetings, Assembly."},
				{token.EOF, ""},
			},
		},
		{
			input: `WHEREAS the customary greeting is "Hello, World!":`,
			tokens: []token.Token{
				{token.WHEREAS, "WHEREAS"},
				{token.STRING, "Hello, World!"},
				{token.EOF, ""},
			},
		},
		{
			input: "WHEREAS an Identifier is capitalized",
			tokens: []token.Token{
				{token.WHEREAS, "WHEREAS"},
				{token.IDENT, "Identifier"},
				{token.EOF, ""},
			},
		},
		{
			input: `WHEREAS the Customary Greeting (hereinafter Greeting) is "Hello, World!":`,
			tokens: []token.Token{
				{token.WHEREAS, "WHEREAS"},
				{token.IDENT, "Customary"},
				{token.IDENT, "Greeting"},
				{token.LPAREN, "("},
				{token.HEREINAFTER, "hereinafter"},
				{token.IDENT, "Greeting"},
				{token.RPAREN, ")"},
				{token.STRING, "Hello, World!"},
				{token.EOF, ""},
			},
		},
		{
			input: `A Resolution Concerning Commentary

WHEREAS a resolution consisting entirely of comments has no effect: now, therefore, 

BE IT RESOLVED that this assembly takes no action.`,
			tokens: []token.Token{
				{token.IDENT, "A"},
				{token.IDENT, "Resolution"},
				{token.IDENT, "Concerning"},
				{token.IDENT, "Commentary"},
				{token.WHEREAS, "WHEREAS"},
				{token.IDENT, "BE"}, // TODO
				{token.IDENT, "IT"}, // TODO
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

func TestScan(t *testing.T) {
	for _, test := range []struct {
		f           func(byte) bool
		input, want string
	}{
		{isLetter, "", ""},
		{isLetter, "a", "a"},
		{isLetter, " ", ""},
		{isLetter, "abc def", "abc"},
		{isLetter, "!abc", ""},
		{isLetter, "abc!", "abc"},

		{isNumeral, "0", "0"},
		{isNumeral, "-3", "-3"},
		{isNumeral, "5", "5"},
		{isNumeral, "256", "256"},
		{isNumeral, "65,536", "65,536"},
		{isNumeral, "24,", "24,"},
		{isNumeral, "24 hours", "24"},
		{isNumeral, " 0", ""},
	} {
		l := New(test.input)
		if got := l.scan(test.f); got != test.want {
			t.Errorf("scan(%v): got %v, want %v", test.input, got, test.want)
		}
	}
}

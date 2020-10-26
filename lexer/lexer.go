package lexer

import (
	"github.com/dkmccandless/assembly/token"
)

// Lexer tokenizes an input string.
type Lexer struct {
	input        string
	pos, readPos int
	ch           byte
}

// New returns a Lexer for input.
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar advances l by one byte and stores the byte at readPos in ch.
// Invariant: While readPos < len(l.input), readPos == pos + 1.
func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
}

// Next returns the next token.Token in l.
func (l *Lexer) Next() token.Token {
	l.skipWhitespace()
	var t token.Token
	switch l.ch {
	case 0:
		return token.Token{token.EOF, ""}
	case '(':
		t = token.Token{token.LPAREN, "("}
	case ')':
		t = token.Token{token.RPAREN, ")"}
	case '-':
		l.readChar()
		if isNumeral(l.ch) {
			return token.Token{token.NUMERAL, "-" + l.scan(isNumeral)}
		}
		return token.Token{token.DASH, "-"}
	default:
		switch {
		case isLetter(l.ch):
			lit := l.scan(isLetter)
			return token.Token{token.Lookup(lit), lit}
		case isDigit(l.ch):
			lit := l.scan(isNumeral)
			return token.Token{token.NUMERAL, lit}
		default:
			t = token.Token{token.COMMENT, string(l.ch)}
		}
	}
	l.readChar()
	return t
}

// scan advances l through all consecutive bytes that satisfy f and returns a string of the bytes read.
func (l *Lexer) scan(f func(b byte) bool) string {
	var s string
	for f(l.ch) {
		s += string(l.ch)
		l.readChar()
	}
	return s
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(b byte) bool { return 'A' <= b && b <= 'Z' || 'a' <= b && b <= 'z' }

// isDigit reports whether b is a digit.
func isDigit(b byte) bool { return '0' <= b && b <= '9' }

// isNumeral reports whether b is a valid character for a numeral literal: a digit, a delimiting comma, or a negative sign.
func isNumeral(b byte) bool { return isDigit(b) || b == ',' || b == '-' }

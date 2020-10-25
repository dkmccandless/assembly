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
	switch l.ch {
	case 0:
		return token.Token{Typ: token.EOF}
	default:
		if isLetter(l.ch) {
			lit := l.scan(isLetter)
			return token.Token{token.Lookup(lit), lit}
		}
	}
	ch := l.ch
	l.readChar()
	return token.Token{token.COMMENT, string(ch)}
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

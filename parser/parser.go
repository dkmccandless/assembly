package parser

import (
	"errors"

	"github.com/dkmccandless/assembly/ast"
	"github.com/dkmccandless/assembly/lexer"
	"github.com/dkmccandless/assembly/token"
)

// Parser parses tokens from a Lexer into an abstract syntax tree.
type Parser struct {
	l *lexer.Lexer

	// cur holds the current token to be parsed.
	cur token.Token

	// peek holds the next token after cur.
	peek token.Token
}

// New returns a pointer to a Parser that parses tokens from l.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.next()
	p.next()
	return p
}

// next consumes the next token from p.l.
func (p *Parser) next() { p.cur, p.peek = p.peek, p.l.Next() }

// curIs reports whether the Type of p.cur is typ.
func (p *Parser) curIs(typ token.Type) bool { return p.cur.Typ == typ }

// peekIs reports whether the Type of p.peek is typ.
func (p *Parser) peekIs(typ token.Type) bool { return p.peek.Typ == typ }

var (
	// Resolution parsing failure errors
	errEarlyResolved = errors.New("no Whereas clause before Resolved clause")
	errLateWhereas   = errors.New("Whereas clause after Resolved clause")
	errNoResolved    = errors.New("no Resolved clause")
	errNoWhereas     = errors.New("no Whereas clause")
)

// ParseResolution parses a Resolution.
// If parsing fails, it returns an error explaining why.
func (p *Parser) ParseResolution() (*ast.Resolution, error) {
	res := &ast.Resolution{}

	// All Whereas clauses must precede all Resolved clauses, and there must be at least one of each.
	var haveWhereas, haveResolved bool

	for !p.curIs(token.EOF) {
		switch p.cur.Typ {
		case token.WHEREAS:
			if haveResolved {
				return nil, errLateWhereas
			}
			haveWhereas = true
		case token.RESOLVED:
			if !haveWhereas {
				for {
					switch p.cur.Typ {
					case token.EOF:
						return nil, errNoWhereas
					case token.WHEREAS:
						return nil, errEarlyResolved
					default:
						p.next()
					}
				}
			}
			haveResolved = true
		}
		p.next()
	}
	if !haveResolved {
		return nil, errNoResolved
	}

	return res, nil
}

// ParseExpr parses an expression.
// If parsing fails, it returns an error explaining why.
func (p *Parser) ParseExpr() (ast.Expr, error) {
	switch {
	case p.cur.IsCardinal():
		return p.parseIntegerLiteral()
	default:
		return nil, errors.New("unrecognized expression")
	}
}

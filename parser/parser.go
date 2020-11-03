package parser

import (
	"errors"
	"fmt"

	"github.com/dkmccandless/assembly/ast"
	"github.com/dkmccandless/assembly/lexer"
	"github.com/dkmccandless/assembly/token"
)

// ErrorList is a list of parsing errors.
// The zero value is an empty ErrorList ready to use.
type ErrorList []error

// Add adds an error to an ErrorList.
func (el ErrorList) Add(err error) { el = append(el, err) }

// Err returns an error equivalent to el, or nil if el is empty.
func (el ErrorList) Err() error {
	if len(el) == 0 {
		return nil
	}
	return el
}

// ErrorList implements the error interface.
func (el ErrorList) Error() string {
	switch len(el) {
	case 0:
		return "no errors"
	case 1:
		return el[0].Error()
	default:
		return fmt.Sprintf("%s (and %v more errors)", el[0], len(el)-1)
	}
}

// Parser parses tokens from a Lexer into an abstract syntax tree.
type Parser struct {
	l      *lexer.Lexer
	errors ErrorList

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

// error adds err to p's ErrorList.
func (p *Parser) error(err error) { p.errors = append(p.errors, err) }

var (
	// Resolution parsing failure errors
	errNoTitle       = errors.New("no title")
	errEarlyResolved = errors.New("no Whereas clause before Resolved clause")
	errLateWhereas   = errors.New("Whereas clause after Resolved clause")
	errNoResolved    = errors.New("no Resolved clause")
	errNoWhereas     = errors.New("no Whereas clause")
)

// ParseResolution parses a Resolution.
// If parsing fails, it returns an error explaining why.
func (p *Parser) ParseResolution() (*ast.Resolution, error) {
	// The Resolution must begin with a title.
	if !p.curIs(token.COMMENT) && !p.curIs(token.IDENT) {
		p.error(errNoTitle)
		return nil, p.errors.Err()
	}

	res := &ast.Resolution{}

	// All Whereas clauses must precede all Resolved clauses, and there must be at least one of each.
	var haveWhereas, haveResolved bool

	for !p.curIs(token.EOF) {
		switch p.cur.Typ {
		case token.WHEREAS:
			if haveResolved {
				p.error(errLateWhereas)
				return nil, p.errors.Err()
			}
			haveWhereas = true
			if stmt := p.parseWhereasStmt(); stmt != nil {
				res.WhereasStmts = append(res.WhereasStmts, stmt)
			}
		case token.RESOLVED:
			if !haveWhereas {
				for !p.curIs(token.EOF) {
					if p.curIs(token.WHEREAS) {
						p.error(errEarlyResolved)
						return nil, p.errors.Err()
					}
					p.next()
				}
				p.error(errNoWhereas)
				return nil, p.errors.Err()
			}
			haveResolved = true
			if stmt := p.parseResolvedStmt(); stmt != nil {
				res.ResolvedStmts = append(res.ResolvedStmts, stmt)
			}
		}
		p.next()
	}
	if !haveResolved {
		p.error(errNoResolved)
		return nil, p.errors.Err()
	}

	return res, p.errors.Err()
}

func (p *Parser) parseWhereasStmt() ast.WhereasStmt {
	for !p.peekIs(token.HEREINAFTER) {
		if p.peekIs(token.WHEREAS) || p.peekIs(token.RESOLVED) || p.peekIs(token.EOF) {
			return nil
		}
		p.next()
	}
	p.next()
	switch p.cur.Typ {
	case token.HEREINAFTER:
		return p.parseDeclStmt()
	default:
		return nil
	}
}

func (p *Parser) parseDeclStmt() *ast.DeclStmt {
	s := &ast.DeclStmt{Token: p.cur}
	p.next()
	for !p.curIs(token.IDENT) {
		p.next()
	}
	s.Name = p.parseIdentifier()
	p.next()
	for !p.cur.IsCardinal() && !p.curIs(token.NUMERAL) && !p.curIs(token.STRING) && !p.curIs(token.IDENT) {
		p.next()
	}
	s.Value = p.ParseExpr()
	return s
}

func (p *Parser) parseResolvedStmt() ast.ResolvedStmt {
	for !p.peekIs(token.PUBLISH) {
		if p.peekIs(token.WHEREAS) || p.peekIs(token.RESOLVED) || p.peekIs(token.EOF) {
			return nil
		}
		p.next()
	}
	p.next()
	switch p.cur.Typ {
	case token.PUBLISH:
		return p.parsePublishStmt()
	default:
		return nil
	}
}

func (p *Parser) parsePublishStmt() *ast.PublishStmt {
	s := &ast.PublishStmt{Token: p.cur}
	p.next()
	for !p.cur.IsCardinal() && !p.curIs(token.STRING) && !p.curIs(token.IDENT) {
		p.next()
	}
	s.Value = p.ParseExpr()
	return s
}

// ParseExpr parses an expression.
func (p *Parser) ParseExpr() ast.Expr {
	switch {
	case p.cur.IsCardinal():
		return p.parseIntegerLiteral()
	case p.curIs(token.NUMERAL):
		// Let parseIntegerLiteral the syntax error
		return p.parseIntegerLiteral()
	case p.curIs(token.STRING):
		return p.parseStringLiteral()
	case p.curIs(token.IDENT):
		return p.parseIdentifier()
	default:
		p.error(errors.New("unrecognized expression"))
		return nil
	}
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	return &ast.Identifier{Token: p.cur, Value: p.cur.Lit}
}

func (p *Parser) parseStringLiteral() *ast.StringLiteral {
	return &ast.StringLiteral{Token: p.cur, Value: p.cur.Lit}
}

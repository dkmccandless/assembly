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
		return nil, errNoTitle
	}

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
			stmt, err := p.parseWhereasStmt()
			if err != nil {
				return nil, err
			}
			if stmt != nil {
				res.WhereasStmts = append(res.WhereasStmts, stmt)
			}
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
			stmt, err := p.parseResolvedStmt()
			if err != nil {
				return nil, err
			}
			if stmt != nil {
				res.ResolvedStmts = append(res.ResolvedStmts, stmt)
			}
		}
		p.next()
	}
	if !haveResolved {
		return nil, errNoResolved
	}

	return res, nil
}

func (p *Parser) parseWhereasStmt() (ast.WhereasStmt, error) {
	for !p.peekIs(token.HEREINAFTER) {
		if p.peekIs(token.WHEREAS) || p.peekIs(token.RESOLVED) || p.peekIs(token.EOF) {
			return nil, nil
		}
		p.next()
	}
	p.next()
	switch p.cur.Typ {
	case token.HEREINAFTER:
		return p.parseDeclStmt()
	default:
		return nil, nil
	}
}

func (p *Parser) parseDeclStmt() (*ast.DeclStmt, error) {
	s := &ast.DeclStmt{Token: p.cur}
	p.next()
	for !p.curIs(token.IDENT) {
		p.next()
	}
	s.Name = p.parseIdentifier()
	p.next()
	for !p.cur.IsCardinal() && !p.curIs(token.STRING) && !p.curIs(token.IDENT) {
		p.next()
	}
	var err error
	s.Value, err = p.ParseExpr()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (p *Parser) parseResolvedStmt() (ast.ResolvedStmt, error) {
	for !p.peekIs(token.PUBLISH) {
		if p.peekIs(token.WHEREAS) || p.peekIs(token.RESOLVED) || p.peekIs(token.EOF) {
			return nil, nil
		}
		p.next()
	}
	p.next()
	switch p.cur.Typ {
	case token.PUBLISH:
		return p.parsePublishStmt()
	default:
		return nil, nil
	}
}

func (p *Parser) parsePublishStmt() (*ast.PublishStmt, error) {
	s := &ast.PublishStmt{Token: p.cur}
	p.next()
	for !p.cur.IsCardinal() && !p.curIs(token.STRING) && !p.curIs(token.IDENT) {
		p.next()
	}
	var err error
	s.Value, err = p.ParseExpr()
	if err != nil {
		return nil, err
	}
	return s, nil
}

// ParseExpr parses an expression.
// If parsing fails, it returns an error explaining why.
func (p *Parser) ParseExpr() (ast.Expr, error) {
	switch {
	case p.cur.IsCardinal():
		return p.parseIntegerLiteral()
	case p.curIs(token.STRING):
		return p.parseStringLiteral(), nil
	case p.curIs(token.IDENT):
		return p.parseIdentifier(), nil
	default:
		return nil, errors.New("unrecognized expression")
	}
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	return &ast.Identifier{Token: p.cur, Value: p.cur.Lit}
}

func (p *Parser) parseStringLiteral() *ast.StringLiteral {
	return &ast.StringLiteral{Token: p.cur, Value: p.cur.Lit}
}

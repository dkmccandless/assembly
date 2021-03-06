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

type usage int

const (
	undeclared usage = iota
	declared
	used
)

type precedence int

const (
	LOWEST precedence = iota
	INFIX
	PREFIX
	POSTFIX
)

// Parser parses tokens from a Lexer into an abstract syntax tree.
type Parser struct {
	l      *lexer.Lexer
	errors ErrorList

	// idents contains all declared identifiers and records whether each has been used.
	idents map[string]usage

	// cur holds the current token to be parsed.
	cur token.Token

	// peek holds the next token after cur.
	peek token.Token
}

// New returns a pointer to a Parser that parses tokens from l.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		idents: make(map[string]usage),
	}
	p.next()
	p.next()
	return p
}

// next consumes the next token from p.l.
func (p *Parser) next() {
	p.cur = p.peek
	var err error
	if p.peek, err = p.l.Next(); err != nil {
		p.error(err)
	}
}

// curIs reports whether the Type of p.cur is typ.
func (p *Parser) curIs(typ token.Type) bool { return p.cur.Typ == typ }

// peekIs reports whether the Type of p.peek is typ.
func (p *Parser) peekIs(typ token.Type) bool { return p.peek.Typ == typ }

func (p *Parser) precedence(t token.Token) precedence {
	if t.IsCardinal() {
		return PREFIX
	}
	switch t.Typ {
	case token.SQUARED, token.CUBED:
		return POSTFIX
	case token.IDENT, token.STRING, token.NUMERAL, token.TWICE, token.THRICE:
		return PREFIX
	case token.LESS:
		return INFIX
	default:
		return LOWEST
	}
}

func (p *Parser) curPrec() precedence  { return p.precedence(p.cur) }
func (p *Parser) peekPrec() precedence { return p.precedence(p.peek) }

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

// redeclaredError indicates the redeclaration of an identifier.
type redeclaredError struct{ ident string }

// redeclaredError implements the error interface.
func (err redeclaredError) Error() string { return fmt.Sprintf("%s redeclared", err.ident) }

// undeclaredError indicates the attempted usage of an undeclared identifier.
type undeclaredError struct{ ident string }

// undeclaredError implements the error interface.
func (err undeclaredError) Error() string { return fmt.Sprintf("%s undeclared", err.ident) }

// unusedError indicates an unused identifier declaration.
type unusedError struct{ ident string }

// unusedError implements the error interface.
func (err unusedError) Error() string { return fmt.Sprintf("%s declared but not used", err.ident) }

// markUsed records that ident has been used.
// If ident was not declared, it records an undeclaredError instead.
func (p *Parser) markUsed(ident string) {
	if p.idents[ident] == undeclared {
		p.error(undeclaredError{ident})
	} else {
		p.idents[ident] = used
	}
}

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
	for id := range p.idents {
		if p.idents[id] != used {
			p.error(unusedError{id})
		}
	}

	return res, p.errors.Err()
}

func (p *Parser) parseWhereasStmt() ast.WhereasStmt {
	for ; !p.peekIs(token.WHEREAS) && !p.peekIs(token.RESOLVED) && !p.peekIs(token.EOF); p.next() {
		if !p.curIs(token.HEREINAFTER) {
			continue
		}
		switch p.cur.Typ {
		case token.HEREINAFTER:
			return p.parseDeclStmt()
		}
	}
	return nil
}

func (p *Parser) parseDeclStmt() *ast.DeclStmt {
	s := &ast.DeclStmt{Token: p.cur}
	p.next()
	for !p.curIs(token.IDENT) {
		p.next()
	}
	s.Name = p.parseIdentifier()
	if id := s.Name.Value; p.idents[id] != undeclared {
		p.error(redeclaredError{id})
	} else {
		p.idents[id] = declared
	}
	p.next()
	for !isExprToken(p.cur) {
		p.next()
	}
	if p.curIs(token.IDENT) {
		p.markUsed(p.cur.Lit)
	}
	s.Value = p.parseExpr(LOWEST)
	return s
}

func (p *Parser) parseResolvedStmt() ast.ResolvedStmt {
	for ; !p.peekIs(token.WHEREAS) && !p.peekIs(token.RESOLVED) && !p.peekIs(token.EOF); p.next() {
		if !p.curIs(token.IDENT) && !p.curIs(token.IF) && !p.curIs(token.PUBLISH) {
			continue
		}
		switch p.cur.Typ {
		case token.IDENT:
			// Assignment if identifier is followed by token.ASSUME
			if !p.peekIs(token.ASSUME) {
				continue
			}
			id := p.parseIdentifier()
			if p.idents[id.Value] == undeclared {
				p.error(undeclaredError{id.Value})
				return nil
			}
			p.next()
			return p.parseAssumeStmt(id)
		case token.IF:
			return p.parseIfStmt()
		case token.PUBLISH:
			return p.parsePublishStmt()
		}
	}
	return nil
}

func (p *Parser) parseAssumeStmt(ident *ast.Identifier) *ast.AssumeStmt {
	s := &ast.AssumeStmt{
		Token: p.cur,
		Name:  ident,
	}
	p.next()
	for !isExprToken(p.cur) {
		p.next()
	}
	s.Value = p.parseExpr(LOWEST)
	return s
}

func (p *Parser) parseIfStmt() *ast.IfStmt {
	s := &ast.IfStmt{Token: p.cur}
	p.next()
	for !isExprToken(p.cur) {
		p.next()
	}
	s.Left = p.parseExpr(LOWEST)
	p.next()
	for !p.curIs(token.EQUALS) && !p.curIs(token.EXCEEDS) {
		p.next()
	}
	s.Relation = p.cur
	p.next()
	for !isExprToken(p.cur) {
		p.next()
	}
	s.Right = p.parseExpr(LOWEST)
	p.next()
	s.Consequence = p.parseResolvedStmt()
	return s
}

func (p *Parser) parsePublishStmt() *ast.PublishStmt {
	s := &ast.PublishStmt{Token: p.cur}
	p.next()
	for !isExprToken(p.cur) {
		p.next()
	}
	s.Value = p.parseExpr(LOWEST)
	return s
}

// parseExpr parses an expression.
func (p *Parser) parseExpr(prec precedence) ast.Expr {
	left := p.parseNullDenotationExpr()
	if left == nil {
		return nil
	}
	// Left-associative
	for prec < p.peekPrec() {
		switch p.peek.Typ {
		case token.LESS:
			p.next()
			left = p.parseInfixExpr(left)
		case token.SQUARED, token.CUBED:
			p.next()
			left = p.parsePostfixExpr(left)
		default:
			return left
		}
	}
	return left
}

// parseNullDenotationExpr parses an expression that begins with a null denotation token
// (representing a literal or prefix operator).
func (p *Parser) parseNullDenotationExpr() ast.Expr {
	if p.cur.IsCardinal() {
		return p.parseIntegerLiteral()
	}
	switch p.cur.Typ {
	case token.IDENT:
		p.markUsed(p.cur.Lit)
		return p.parseIdentifier()
	case token.STRING:
		return p.parseStringLiteral()
	case token.NUMERAL:
		// Let parseIntegerLiteral record the syntax error
		return p.parseIntegerLiteral()
	case token.TWICE, token.THRICE:
		return p.parseUnaryPrefixExpr()
	case token.SUM, token.PRODUCT, token.QUOTIENT, token.REMAINDER:
		return p.parseBinaryPrefixExpr()
	default:
		p.error(fmt.Errorf("unrecognized expression %v", p.cur.Lit))
		return nil
	}
}

func (p *Parser) parseUnaryPrefixExpr() ast.Expr {
	expr := &ast.UnaryPrefixExpr{Token: p.cur}
	p.next()
	expr.Right = p.parseExpr(PREFIX)
	return expr
}

func (p *Parser) parseBinaryPrefixExpr() ast.Expr {
	expr := &ast.BinaryPrefixExpr{Token: p.cur}
	p.next()
	expr.First = p.parseExpr(LOWEST)
	p.next()
	expr.Second = p.parseExpr(PREFIX)
	return expr
}

// parsePostfixExpr parses a postfix expression: an expression in left denotation context
// that does not accept a following expression.
func (p *Parser) parsePostfixExpr(left ast.Expr) ast.Expr {
	expr := &ast.PostfixExpr{
		Token: p.cur,
		Left:  left,
	}
	return expr
}

// parseInfixExpr parses an infix expression: an expression in left denotation context
// that expects a following expression.
func (p *Parser) parseInfixExpr(left ast.Expr) ast.Expr {
	expr := &ast.InfixExpr{
		Token: p.cur,
		Left:  left,
	}
	prec := p.curPrec()
	p.next()
	expr.Right = p.parseExpr(prec)
	return expr
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	return &ast.Identifier{Token: p.cur, Value: p.cur.Lit}
}

func (p *Parser) parseStringLiteral() *ast.StringLiteral {
	return &ast.StringLiteral{Token: p.cur, Value: p.cur.Lit}
}

// isExprToken reports whether t can begin an ast.Expr.
func isExprToken(t token.Token) bool {
	switch t.Typ {
	case token.STRING, token.IDENT,
		token.TWICE, token.THRICE,
		token.SUM, token.PRODUCT, token.QUOTIENT, token.REMAINDER:
		return true
	case token.NUMERAL:
		// Let parseIntegerLiteral record the syntax error
		return true
	default:
		return t.IsCardinal()
	}
}

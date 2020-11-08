package ast

import (
	"fmt"
	"strconv"

	"github.com/dkmccandless/assembly/token"
)

// All node types implement the Node interface.
type Node interface {
	String() string
}

// A Resolution represents an Assembly program.
type Resolution struct {
	WhereasStmts  []WhereasStmt
	ResolvedStmts []ResolvedStmt
}

func (r *Resolution) String() string { return "" }

// Statements that can occur in Whereas clauses implement the WhereasStmt interface.
type WhereasStmt interface {
	Node
	whStmtNode()
}

type DeclStmt struct {
	Token token.Token // token.HEREINAFTER
	Name  *Identifier
	Value Expr
}

func (s *DeclStmt) whStmtNode()    {}
func (s *DeclStmt) String() string { return s.Token.Lit }

// Statements that can occur in Resolved clauses implement the ResolvedStmt interface.
type ResolvedStmt interface {
	Node
	resStmtNode()
}

type PublishStmt struct {
	Token token.Token // token.PUBLISH
	Value Expr
}

func (s *PublishStmt) resStmtNode()   {}
func (s *PublishStmt) String() string { return s.Token.Lit }

// Expressions implement the Expr interface.
type Expr interface {
	Node
	exprNode()
}

type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

func (e *Identifier) exprNode()      {}
func (e *Identifier) String() string { return e.Value }

type IntegerLiteral struct {
	Token token.Token // token.INTEGER
	Value int64
}

func (e *IntegerLiteral) exprNode()      {}
func (e *IntegerLiteral) String() string { return strconv.Itoa(int(e.Value)) }

type StringLiteral struct {
	Token token.Token // token.STRING
	Value string
}

func (e *StringLiteral) exprNode()      {}
func (e *StringLiteral) String() string { return e.Value }

type InfixExpr struct {
	Token       token.Token // e.g. token.LESS
	Left, Right Expr
}

func (e *InfixExpr) exprNode()      {}
func (e *InfixExpr) String() string { return fmt.Sprintf("%v %v %v", e.Left, e.Token.Lit, e.Right) }

type UnaryPrefixExpr struct {
	Token token.Token // e.g. token.THRICE
	Right Expr
}

func (e *UnaryPrefixExpr) exprNode()      {}
func (e *UnaryPrefixExpr) String() string { return fmt.Sprintf("%v %v", e.Token.Lit, e.Right) }

type BinaryPrefixExpr struct {
	Token         token.Token // e.g. token.SUM
	First, Second Expr
}

func (e *BinaryPrefixExpr) exprNode() {}
func (e *BinaryPrefixExpr) String() string {
	return fmt.Sprintf("%v %v %v", e.Token.Lit, e.First, e.Second)
}

package ast

import (
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

// Statements that can occur in Whereas clauses implement the WhereasStmt interface.
type WhereasStmt interface {
	Node
	whStmtNode()
}

// Statements that can occur in Resolved clauses implement the ResolvedStmt interface.
type ResolvedStmt interface {
	Node
	resStmtNode()
}

// Expressions implement the Expr interface.
type Expr interface {
	Node
	exprNode()
}

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

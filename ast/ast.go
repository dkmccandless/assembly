package ast

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

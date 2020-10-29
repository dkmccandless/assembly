package eval

import (
	"github.com/dkmccandless/assembly/ast"
	"github.com/dkmccandless/assembly/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	default:
		return nil
	}
}

package eval

import (
	"github.com/dkmccandless/assembly/ast"
	"github.com/dkmccandless/assembly/object"
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.Identifier:
		// An ast.Identifier is created for every capitalized non-keyword;
		// return nil if the "identifier" is not in env.
		if obj, ok := env.Get(node.Value); ok {
			return obj
		}
		return nil
	default:
		return nil
	}
}

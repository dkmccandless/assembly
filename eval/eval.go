package eval

import (
	"fmt"

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
	case *ast.DeclStmt:
		if val := Eval(node.Value, env); val != nil {
			env.Set(node.Name.Value, val)
		}
	case *ast.PublishStmt:
		if val := Eval(node.Value, env); val != nil {
			fmt.Println(val.Inspect())
		}
	}
	return nil
}

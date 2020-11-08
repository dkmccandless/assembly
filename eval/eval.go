package eval

import (
	"fmt"

	"github.com/dkmccandless/assembly/ast"
	"github.com/dkmccandless/assembly/object"
	"github.com/dkmccandless/assembly/token"
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.UnaryPrefixExpr:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalUnaryPrefixExpr(node.Token, right)
	case *ast.BinaryPrefixExpr:
		first := Eval(node.First, env)
		if isError(first) {
			return first
		}
		second := Eval(node.Second, env)
		if isError(second) {
			return second
		}
		return evalBinaryPrefixExpr(node.Token, first, second)
	case *ast.InfixExpr:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpr(node.Token, left, right)
	case *ast.PostfixExpr:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		return evalPostfixExpr(node.Token, left)
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
	case *ast.AssumeStmt:
		if val := Eval(node.Value, env); val != nil {
			env.Set(node.Name.Value, val)
		}
	case *ast.PublishStmt:
		if val := Eval(node.Value, env); val != nil {
			fmt.Println(val.Inspect())
		}
	case *ast.Resolution:
		for _, wh := range node.WhereasStmts {
			if err := Eval(wh, env); err != nil {
				return err
			}
		}
		for _, res := range node.ResolvedStmts {
			if err := Eval(res, env); err != nil {
				return err
			}
		}
	}
	return nil
}

func evalUnaryPrefixExpr(t token.Token, right object.Object) object.Object {
	if right.Type() != object.INTEGER {
		return nonNumericError(right)
	}
	r := right.(*object.Integer).Value
	switch t.Typ {
	case token.TWICE:
		return &object.Integer{2 * r}
	case token.THRICE:
		return &object.Integer{3 * r}
	default:
		return &object.Error{fmt.Sprintf("unknown operator %v %v", t.Lit, r)}
	}
}

func evalBinaryPrefixExpr(t token.Token, first, second object.Object) object.Object {
	if first.Type() != object.INTEGER {
		return nonNumericError(first)
	}
	a := first.(*object.Integer).Value
	if second.Type() != object.INTEGER {
		return nonNumericError(second)
	}
	b := second.(*object.Integer).Value
	switch t.Typ {
	case token.SUM:
		return &object.Integer{a + b}
	case token.PRODUCT:
		return &object.Integer{a * b}
	case token.QUOTIENT:
		return &object.Integer{a / b}
	case token.REMAINDER:
		return &object.Integer{a % b}
	default:
		return &object.Error{fmt.Sprintf("unknown operator %v %v %v", t.Lit, a, b)}
	}
}

func evalInfixExpr(t token.Token, left, right object.Object) object.Object {
	if left.Type() != object.INTEGER {
		return nonNumericError(left)
	}
	l := left.(*object.Integer).Value
	if right.Type() != object.INTEGER {
		return nonNumericError(right)
	}
	r := right.(*object.Integer).Value
	switch t.Typ {
	case token.LESS:
		return &object.Integer{l - r}
	default:
		return &object.Error{fmt.Sprintf("unknown operator %v %v %v", l, t.Lit, r)}
	}
}

func evalPostfixExpr(t token.Token, left object.Object) object.Object {
	if left.Type() != object.INTEGER {
		return nonNumericError(left)
	}
	l := left.(*object.Integer).Value
	switch t.Typ {
	case token.SQUARED:
		return &object.Integer{l * l}
	case token.CUBED:
		return &object.Integer{l * l * l}
	default:
		return &object.Error{fmt.Sprintf("unknown operator %v %v", l, t.Lit)}
	}
}

// nonNumericError records that obj occurs in an expression context that requires a numeric form.
func nonNumericError(obj object.Object) *object.Error {
	return &object.Error{fmt.Sprintf("non-numeric %s in numeric context", obj.Inspect())}
}

func isError(obj object.Object) bool { return obj != nil && obj.Type() == object.ERROR }

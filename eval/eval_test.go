package eval

import (
	"math"
	"reflect"
	"testing"

	"github.com/dkmccandless/assembly/ast"
	"github.com/dkmccandless/assembly/object"
	"github.com/dkmccandless/assembly/token"
)

func TestEvalIntegerExpr(t *testing.T) {
	for _, test := range []struct {
		ast ast.Expr
		obj object.Object
	}{
		{
			&ast.IntegerLiteral{token.Token{token.INTEGER, "0"}, 0},
			&object.Integer{0},
		},
		{
			&ast.IntegerLiteral{token.Token{token.INTEGER, "1"}, 1},
			&object.Integer{1},
		},
		{
			&ast.IntegerLiteral{token.Token{token.INTEGER, "-3000000000000"}, -3000000000000},
			&object.Integer{-3000000000000},
		},
		{
			&ast.IntegerLiteral{token.Token{token.INTEGER, "-9223372036854775808"}, math.MinInt64},
			&object.Integer{math.MinInt64},
		},
		{
			&ast.IntegerLiteral{token.Token{token.INTEGER, "9223372036854775807"}, math.MaxInt64},
			&object.Integer{math.MaxInt64},
		},
	} {
		if obj := Eval(test.ast, object.NewEnvironment()); !reflect.DeepEqual(obj, test.obj) {
			t.Errorf("EvalIntegerExpr(%+v): got %+v, want %+v", test.ast, obj, test.obj)
		}
	}
}

func TestEvalStringExpr(t *testing.T) {
	for _, test := range []struct {
		ast ast.Expr
		obj object.Object
	}{
		{
			&ast.StringLiteral{token.Token{token.STRING, ""}, ""},
			&object.String{""},
		},
		{
			&ast.StringLiteral{token.Token{token.STRING, "WHEREAS"}, "WHEREAS"},
			&object.String{"WHEREAS"},
		},
		{
			&ast.StringLiteral{token.Token{token.STRING, "zero (0)"}, "zero (0)"},
			&object.String{"zero (0)"},
		},
		{
			&ast.StringLiteral{token.Token{token.STRING, "Greetings, Assembly."}, "Greetings, Assembly."},
			&object.String{"Greetings, Assembly."},
		},
	} {
		if obj := Eval(test.ast, object.NewEnvironment()); !reflect.DeepEqual(obj, test.obj) {
			t.Errorf("EvalStringExpr(%+v): got %+v, want %+v", test.ast, obj, test.obj)
		}
	}
}

func TestEvalIdentifier(t *testing.T) {
	for _, test := range []struct {
		ast ast.Expr
		obj object.Object
	}{
		{
			&ast.Identifier{token.Token{token.IDENT, "Greeting"}, "Greeting"},
			&object.String{"Greeting ok"},
		},
		{
			&ast.Identifier{token.Token{token.IDENT, "Quantity"}, "Quantity"},
			&object.String{"Quantity ok"},
		},
		{
			&ast.Identifier{token.Token{token.IDENT, "Answer"}, "Answer"},
			&object.String{"Answer ok"},
		},
	} {
		env := object.NewEnvironment()
		if obj := Eval(test.ast, env); obj != nil {
			t.Errorf("EvalIdentifier(%v): got %T (%+v) before Set", test, obj, obj)
		}
		id := test.ast.(*ast.Identifier).Value
		want := &object.String{Value: id + " ok"}
		env.Set(id, want)
		if obj, ok := Eval(test.ast, env).(*object.String); !ok {
			t.Errorf("EvalIdentifier(%v): got %T (%+v) after Set, want %T (%+v)", test, obj, obj, want, want)
		}
	}
}

func TestEvalDeclStmt(t *testing.T) {
	for _, ast := range []*ast.DeclStmt{
		&ast.DeclStmt{
			Token: token.Token{Typ: token.HEREINAFTER, Lit: "hereinafter"},
			Name: &ast.Identifier{
				Token: token.Token{Typ: token.IDENT, Lit: "Greeting"},
				Value: "Greeting",
			},
			Value: &ast.StringLiteral{
				Token: token.Token{Typ: token.STRING, Lit: "Hello, World!"},
				Value: "Hello, World!",
			},
		},
		&ast.DeclStmt{
			Token: token.Token{Typ: token.HEREINAFTER, Lit: "hereinafter"},
			Name: &ast.Identifier{
				Token: token.Token{Typ: token.IDENT, Lit: "Answer"},
				Value: "Answer",
			},
			Value: &ast.IntegerLiteral{
				Token: token.Token{Typ: token.INTEGER, Lit: "42"},
				Value: 42,
			},
		},
	} {
		env := object.NewEnvironment()
		if obj := Eval(ast.Name, env); obj != nil {
			t.Errorf("EvalDeclStmt(Identifier %+v): got %T (%+v) before declaration", ast.Name, obj, obj)
		}
		if obj := Eval(ast, env); obj != nil {
			t.Errorf("EvalDeclStmt(%+v): got %T (%+v) from declaration", ast, obj, obj)
		}
		if obj := Eval(ast.Name, env); obj == nil {
			t.Errorf("EvalDeclStmt(Identifier %+v): got %T (%+v) after declaration", ast.Name, obj, obj)
		}
	}
}

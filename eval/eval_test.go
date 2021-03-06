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
	for _, test := range []struct {
		stmt *ast.DeclStmt
		obj  object.Object
	}{
		{
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
			&object.String{"Hello, World!"},
		},
		{
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
			&object.Integer{42},
		},
		{
			&ast.DeclStmt{
				Token: token.Token{Typ: token.HEREINAFTER, Lit: "hereinafter"},
				Name: &ast.Identifier{
					Token: token.Token{Typ: token.IDENT, Lit: "Dozen"},
					Value: "Dozen",
				},
				Value: &ast.BinaryPrefixExpr{
					Token:  token.Token{Typ: token.SUM, Lit: "sum"},
					First:  &ast.IntegerLiteral{token.Token{token.INTEGER, "10"}, 10},
					Second: &ast.IntegerLiteral{token.Token{token.INTEGER, "2"}, 2},
				},
			},
			&object.Integer{12},
		},
	} {
		env := object.NewEnvironment()
		id := test.stmt.Name
		if obj := Eval(id, env); obj != nil {
			t.Errorf("EvalDeclStmt(Identifier %+v): got %T (%+v) before declaration", id, obj, obj)
		}
		if obj := Eval(test.stmt, env); obj != nil {
			t.Errorf("EvalDeclStmt(%+v): got %T (%+v) from declaration", test.stmt, obj, obj)
		}
		if obj := Eval(id, env); !reflect.DeepEqual(obj, test.obj) {
			t.Errorf("EvalDeclStmt(Identifier %+v): got %T (%+v) after declaration", id, obj, obj)
		}
	}
}

func TestAssumeStmt(t *testing.T) {
	for _, test := range []struct {
		stmt *ast.AssumeStmt
		obj  object.Object
	}{
		{
			&ast.AssumeStmt{
				Token: token.Token{Typ: token.ASSUME, Lit: "assume"},
				Name: &ast.Identifier{
					Token: token.Token{Typ: token.IDENT, Lit: "Greeting"},
					Value: "Greeting",
				},
				Value: &ast.StringLiteral{
					Token: token.Token{Typ: token.STRING, Lit: "Greetings, Assembly."},
					Value: "Greetings, Assembly.",
				},
			},
			&object.String{"Greetings, Assembly."},
		},
		{
			&ast.AssumeStmt{
				Token: token.Token{Typ: token.ASSUME, Lit: "assume"},
				Name: &ast.Identifier{
					Token: token.Token{Typ: token.IDENT, Lit: "Total"},
					Value: "Total",
				},
				Value: &ast.IntegerLiteral{
					Token: token.Token{Typ: token.INTEGER, Lit: "3"},
					Value: 3,
				},
			},
			&object.Integer{3},
		},
		{
			&ast.AssumeStmt{
				Token: token.Token{Typ: token.ASSUME, Lit: "assume"},
				Name: &ast.Identifier{
					Token: token.Token{Typ: token.IDENT, Lit: "Total"},
					Value: "Total",
				},
				Value: &ast.BinaryPrefixExpr{
					Token: token.Token{Typ: token.SUM, Lit: "sum"},
					First: &ast.Identifier{
						Token: token.Token{Typ: token.IDENT, Lit: "Total"},
						Value: "Total",
					},
					Second: &ast.IntegerLiteral{
						Token: token.Token{Typ: token.INTEGER, Lit: "1"},
						Value: 1,
					},
				},
			},
			&object.Integer{11},
		},
	} {
		env := object.NewEnvironment()
		env.Set("Greeting", &object.String{"Hello, World!"})
		env.Set("Total", &object.Integer{10})
		id := test.stmt.Name.Value
		err := Eval(test.stmt, env)
		obj, ok := env.Get(id)
		if err != nil || !reflect.DeepEqual(obj, test.obj) || !ok {
			t.Errorf("EvalAssumeStmt(%v): got %T (%+v), want %T (%+v)", id, obj, obj, test.obj, test.obj)
		}
	}
}

func TestIfStmt(t *testing.T) {
	for _, test := range []struct {
		stmt *ast.IfStmt
		id   string
		obj  object.Object
	}{
		{
			&ast.IfStmt{
				Token:    token.Token{token.IF, "if"},
				Left:     &ast.Identifier{token.Token{token.IDENT, "Error"}, "Error"},
				Right:    &ast.IntegerLiteral{token.Token{token.INTEGER, "1"}, 1},
				Relation: token.Token{token.EQUALS, "equals"},
				Consequence: &ast.AssumeStmt{
					Token: token.Token{token.ASSUME, "assume"},
					Name:  &ast.Identifier{token.Token{token.IDENT, "Error"}, "Error"},
					Value: &ast.BinaryPrefixExpr{
						Token:  token.Token{token.SUM, "sum"},
						First:  &ast.Identifier{token.Token{token.IDENT, "Error"}, "Error"},
						Second: &ast.IntegerLiteral{token.Token{token.INTEGER, "1"}, 1},
					},
				},
			},
			"Error",
			&object.Integer{1},
		},
		{
			&ast.IfStmt{
				Token:    token.Token{token.IF, "if"},
				Left:     &ast.Identifier{token.Token{token.IDENT, "Quorum"}, "Quorum"},
				Right:    &ast.Identifier{token.Token{token.IDENT, "Attendance"}, "Attendance"},
				Relation: token.Token{token.EXCEEDS, "exceeds"},
				Consequence: &ast.AssumeStmt{
					Token: token.Token{token.ASSUME, "assume"},
					Name:  &ast.Identifier{token.Token{token.IDENT, "Message"}, "Message"},
					Value: &ast.StringLiteral{token.Token{token.STRING, "This Assembly lacks a quorum."}, "This Assembly lacks a quorum."},
				},
			},
			"Message",
			&object.String{"This Assembly lacks a quorum."},
		},
	} {
		env := object.NewEnvironment()
		env.Set("Error", &object.Integer{0})
		env.Set("Quorum", &object.Integer{10})
		env.Set("Attendance", &object.Integer{12})
		env.Set("Message", &object.String{"No messages."})
		if err := Eval(test.stmt, env); err != nil {
			t.Errorf("EvalIfStmt(%+v): got error %v", test.stmt, err)
		}
		obj, _ := env.Get(test.id)
		if reflect.DeepEqual(obj, test.obj) {
			t.Errorf("EvalIfStmt(%+v): consequence evaluated on false condition", test.stmt)
		}
		env.Set("Error", &object.Integer{1})
		env.Set("Attendance", &object.Integer{8})
		if err := Eval(test.stmt, env); err != nil {
			t.Errorf("EvalIfStmt(%+v): got error %v", test.stmt, err)
		}
		obj, _ = env.Get(test.id)
		if !reflect.DeepEqual(obj, test.obj) {
			t.Errorf("EvaIfStmt(%+v, %v): consequence not evaluated on true condition", test.stmt, test.id)
		}
	}
}

func TestEvalUnaryPrefixExpr(t *testing.T) {
	for _, test := range []struct {
		ast ast.Expr
		obj object.Object
	}{
		{
			&ast.UnaryPrefixExpr{
				Token: token.Token{token.TWICE, "twice"},
				Right: &ast.IntegerLiteral{token.Token{token.INTEGER, "3"}, 3},
			},
			&object.Integer{6},
		},
		{
			&ast.UnaryPrefixExpr{
				Token: token.Token{token.THRICE, "thrice"},
				Right: &ast.IntegerLiteral{token.Token{token.INTEGER, "4"}, 4},
			},
			&object.Integer{12},
		},
		{
			&ast.UnaryPrefixExpr{
				Token: token.Token{token.THRICE, "thrice"},
				Right: &ast.UnaryPrefixExpr{
					Token: token.Token{token.TWICE, "twice"},
					Right: &ast.IntegerLiteral{token.Token{token.INTEGER, "-1"}, -1},
				},
			},
			&object.Integer{-6},
		},
	} {
		if obj := Eval(test.ast, object.NewEnvironment()); !reflect.DeepEqual(obj, test.obj) {
			t.Errorf("EvalUnaryPrefixExpr(%+v): got %+v, want %+v", test.ast, obj, test.obj)
		}
	}
}

func TestEvalBinaryPrefixExpr(t *testing.T) {
	for _, test := range []struct {
		ast ast.Expr
		obj object.Object
	}{
		{
			&ast.BinaryPrefixExpr{
				Token:  token.Token{token.SUM, "sum"},
				First:  &ast.IntegerLiteral{token.Token{token.INTEGER, "1"}, 1},
				Second: &ast.IntegerLiteral{token.Token{token.INTEGER, "1"}, 1},
			},
			&object.Integer{2},
		},
		{
			&ast.BinaryPrefixExpr{
				Token:  token.Token{token.PRODUCT, "product"},
				First:  &ast.IntegerLiteral{token.Token{token.INTEGER, "2"}, 2},
				Second: &ast.IntegerLiteral{token.Token{token.INTEGER, "3"}, 3},
			},
			&object.Integer{6},
		},
		{
			&ast.BinaryPrefixExpr{
				Token:  token.Token{token.QUOTIENT, "quotient"},
				First:  &ast.IntegerLiteral{token.Token{token.INTEGER, "17"}, 17},
				Second: &ast.IntegerLiteral{token.Token{token.INTEGER, "5"}, 5},
			},
			&object.Integer{3},
		},
		{
			&ast.BinaryPrefixExpr{
				Token:  token.Token{token.REMAINDER, "remainder"},
				First:  &ast.IntegerLiteral{token.Token{token.INTEGER, "17"}, 17},
				Second: &ast.IntegerLiteral{token.Token{token.INTEGER, "5"}, 5},
			},
			&object.Integer{2},
		},
	} {
		if obj := Eval(test.ast, object.NewEnvironment()); !reflect.DeepEqual(obj, test.obj) {
			t.Errorf("EvalBinaryPrefixExpr(%+v): got %+v, want %+v", test.ast, obj, test.obj)
		}
	}
}

func TestEvalInfixExpr(t *testing.T) {
	for _, test := range []struct {
		ast ast.Expr
		obj object.Object
	}{
		{
			&ast.InfixExpr{
				Token: token.Token{token.LESS, "less"},
				Left:  &ast.IntegerLiteral{token.Token{token.INTEGER, "3"}, 3},
				Right: &ast.IntegerLiteral{token.Token{token.INTEGER, "2"}, 2},
			},
			&object.Integer{1},
		},
	} {
		if obj := Eval(test.ast, object.NewEnvironment()); !reflect.DeepEqual(obj, test.obj) {
			t.Errorf("EvalInfixExpr(%+v): got %+v, want %+v", test.ast, obj, test.obj)
		}
	}
}

func TestEvalPostfixExpr(t *testing.T) {
	for _, test := range []struct {
		ast ast.Expr
		obj object.Object
	}{
		{
			&ast.PostfixExpr{
				Token: token.Token{token.SQUARED, "squared"},
				Left:  &ast.IntegerLiteral{token.Token{token.INTEGER, "3"}, 3},
			},
			&object.Integer{9},
		},
		{
			&ast.PostfixExpr{
				Token: token.Token{token.CUBED, "cubed"},
				Left:  &ast.IntegerLiteral{token.Token{token.INTEGER, "4"}, 4},
			},
			&object.Integer{64},
		},
		{
			&ast.PostfixExpr{
				Token: token.Token{token.SQUARED, "squared"},
				Left: &ast.PostfixExpr{
					Token: token.Token{token.CUBED, "cubed"},
					Left:  &ast.IntegerLiteral{token.Token{token.INTEGER, "10"}, 10},
				},
			},
			&object.Integer{1e6},
		},
	} {
		if obj := Eval(test.ast, object.NewEnvironment()); !reflect.DeepEqual(obj, test.obj) {
			t.Errorf("EvalPostfixExpr(%+v): got %+v, want %+v", test.ast, obj, test.obj)
		}
	}
}

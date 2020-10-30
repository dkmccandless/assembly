package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/dkmccandless/assembly/ast"
	"github.com/dkmccandless/assembly/lexer"
	"github.com/dkmccandless/assembly/token"
)

func TestParseResolution(t *testing.T) {
	for _, test := range []struct {
		input string
		ast   *ast.Resolution
		err   error
	}{
		{"whereas", nil, errNoTitle},
		{"title whereas", nil, errNoResolved},
		{"title resolved", nil, errNoWhereas},
		{"title resolved whereas resolved", nil, errEarlyResolved},
		{"title whereas resolved whereas", nil, errLateWhereas},
		{"title whereas resolved", &ast.Resolution{}, nil},
		{"title whereas whereas resolved", &ast.Resolution{}, nil},
		{"title whereas resolved resolved", &ast.Resolution{}, nil},
		{"title whereas whereas resolved resolved", &ast.Resolution{}, nil},
	} {
		p := New(lexer.New(test.input))
		if ast, err := p.ParseResolution(); !reflect.DeepEqual(ast, test.ast) || err != test.err {
			t.Errorf("ParseResolution(%v): got %v, %v; want %v, %v", test.input, ast, err, test.ast, test.err)
		}
	}
}

var stringTests = []string{
	"",
	"WHEREAS",
	"zero (0)",
	"Greetings, Assembly.",
}

func TestParseStringLiteral(t *testing.T) {
	for _, test := range stringTests {
		input := fmt.Sprintf("\"%v\"", test)
		want := &ast.StringLiteral{
			Token: token.Token{Typ: token.STRING, Lit: test},
			Value: test,
		}
		p := New(lexer.New(input))
		if got := p.parseStringLiteral(); !reflect.DeepEqual(got, want) {
			t.Errorf("parseStringLiteral(%v): got %#v, want %#v", input, got, want)
		}
	}
}

func TestParseExpr(t *testing.T) {
	for _, test := range integerTests {
		input := fmt.Sprintf("%v (%v)", test.car, test.num)
		p := New(lexer.New(input))
		expr, err := p.ParseExpr()
		if err != nil {
			t.Errorf("ParseExpr(%v): got error %v", input, err)
		}
		e, ok := expr.(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("ParseExpr(%v): got %T (%+v)", input, expr, expr)
		}
		if e.Value != test.n {
			t.Errorf("ParseExpr(%v): got %v", input, e.Value)
		}
	}
	for _, test := range stringTests {
		input := fmt.Sprintf("\"%v\"", test)
		p := New(lexer.New(input))
		expr, err := p.ParseExpr()
		if err != nil {
			t.Errorf("ParseExpr(%v): got error %v", input, err)
		}
		e, ok := expr.(*ast.StringLiteral)
		if !ok {
			t.Errorf("ParseExpr(%v): got %T (%+v)", input, expr, expr)
		}
		if e.Value != test {
			t.Errorf("ParseExpr(%v): got %v", input, e.Value)
		}
	}
}

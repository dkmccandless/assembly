package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/dkmccandless/assembly/ast"
	"github.com/dkmccandless/assembly/lexer"
)

func TestParseResolution(t *testing.T) {
	for _, test := range []struct {
		input string
		ast   *ast.Resolution
		err   error
	}{
		{"whereas", nil, errNoResolved},
		{"resolved", nil, errNoWhereas},
		{"resolved whereas resolved", nil, errEarlyResolved},
		{"whereas resolved whereas", nil, errLateWhereas},
		{"whereas resolved", &ast.Resolution{}, nil},
	} {
		p := New(lexer.New(test.input))
		if ast, err := p.ParseResolution(); !reflect.DeepEqual(ast, test.ast) || err != test.err {
			t.Errorf("ParseResolution(%v): got %v, %v; want %v, %v", test.input, ast, err, test.ast, test.err)
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
		i, ok := expr.(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("ParseExpr(%v): got %T (%+v)", input, expr, expr)
		}
		if i.Value != test.n {
			t.Errorf("ParseExpr(%v): got %v", input, i.Value)
		}
	}
}

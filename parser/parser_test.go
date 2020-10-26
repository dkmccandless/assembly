package parser

import (
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

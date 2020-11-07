package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/dkmccandless/assembly/ast"
	"github.com/dkmccandless/assembly/lexer"
	"github.com/dkmccandless/assembly/token"
)

// lastError returns the last error in p.errors, or nil if no errors were recorded.
func (p *Parser) lastError() error {
	if len(p.errors) != 0 {
		return p.errors[len(p.errors)-1]
	}
	return nil
}

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
		{
			`title whereas resolved publish "Hello, World!"`,
			&ast.Resolution{
				ResolvedStmts: []ast.ResolvedStmt{
					&ast.PublishStmt{
						Token: token.Token{Typ: token.PUBLISH, Lit: "publish"},
						Value: &ast.StringLiteral{
							Token: token.Token{Typ: token.STRING, Lit: "Hello, World!"},
							Value: "Hello, World!",
						},
					},
				},
			},
			nil,
		},
		{
			`title whereas the Customary Greeting is "Hello, World!" resolved publish Greeting`,
			&ast.Resolution{
				ResolvedStmts: []ast.ResolvedStmt{
					&ast.PublishStmt{
						Token: token.Token{Typ: token.PUBLISH, Lit: "publish"},
						Value: &ast.Identifier{
							Token: token.Token{Typ: token.IDENT, Lit: "Greeting"},
							Value: "Greeting",
						},
					},
				},
			},
			undeclaredError{"Greeting"},
		},
		{
			`title whereas the Customary Greeting (hereinafter Greeting) is "Hello, World!" resolved`,
			&ast.Resolution{
				WhereasStmts: []ast.WhereasStmt{
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
				},
			},
			unusedError{"Greeting"},
		},
		{
			`title whereas the Customary Greeting (hereinafter Greeting) is "Hello, World!" resolved publish "Hello, World!`,
			&ast.Resolution{
				WhereasStmts: []ast.WhereasStmt{
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
				},
				ResolvedStmts: []ast.ResolvedStmt{
					&ast.PublishStmt{
						Token: token.Token{Typ: token.PUBLISH, Lit: "publish"},
						Value: &ast.StringLiteral{
							Token: token.Token{Typ: token.STRING, Lit: "Hello, World!"},
							Value: "Hello, World!",
						},
					},
				},
			},
			unusedError{"Greeting"},
		},
		{
			`title whereas the Customary Greeting (hereinafter Greeting) is "Hello, World!" whereas the Customary Greeting (hereinafter Greeting) is "Hello, World!" resolved publish Greeting`,
			&ast.Resolution{
				WhereasStmts: []ast.WhereasStmt{
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
							Token: token.Token{Typ: token.IDENT, Lit: "Greeting"},
							Value: "Greeting",
						},
						Value: &ast.StringLiteral{
							Token: token.Token{Typ: token.STRING, Lit: "Hello, World!"},
							Value: "Hello, World!",
						},
					},
				},
				ResolvedStmts: []ast.ResolvedStmt{
					&ast.PublishStmt{
						Token: token.Token{Typ: token.PUBLISH, Lit: "publish"},
						Value: &ast.Identifier{
							Token: token.Token{Typ: token.IDENT, Lit: "Greeting"},
							Value: "Greeting",
						},
					},
				},
			},
			redeclaredError{"Greeting"},
		},
		{
			`title whereas the Customary Greeting (hereinafter Greeting) is "Hello, World!" resolved publish Greeting`,
			&ast.Resolution{
				WhereasStmts: []ast.WhereasStmt{
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
				},
				ResolvedStmts: []ast.ResolvedStmt{
					&ast.PublishStmt{
						Token: token.Token{Typ: token.PUBLISH, Lit: "publish"},
						Value: &ast.Identifier{
							Token: token.Token{Typ: token.IDENT, Lit: "Greeting"},
							Value: "Greeting",
						},
					},
				},
			},
			nil,
		},
		{
			`title
whereas the Amount in Stock (hereinafter Stock) is ninety-nine (99)
whereas the Current Quantity (hereinafter Quantity) is Stock
resolved publish Quantity`,
			&ast.Resolution{
				WhereasStmts: []ast.WhereasStmt{
					&ast.DeclStmt{
						Token: token.Token{Typ: token.HEREINAFTER, Lit: "hereinafter"},
						Name: &ast.Identifier{
							Token: token.Token{Typ: token.IDENT, Lit: "Stock"},
							Value: "Stock",
						},
						Value: &ast.IntegerLiteral{
							Token: token.Token{Typ: token.INTEGER, Lit: "99"},
							Value: 99,
						},
					},
					&ast.DeclStmt{
						Token: token.Token{Typ: token.HEREINAFTER, Lit: "hereinafter"},
						Name: &ast.Identifier{
							Token: token.Token{Typ: token.IDENT, Lit: "Quantity"},
							Value: "Quantity",
						},
						Value: &ast.Identifier{
							Token: token.Token{Typ: token.IDENT, Lit: "Stock"},
							Value: "Stock",
						},
					},
				},
				ResolvedStmts: []ast.ResolvedStmt{
					&ast.PublishStmt{
						Token: token.Token{Typ: token.PUBLISH, Lit: "publish"},
						Value: &ast.Identifier{
							Token: token.Token{Typ: token.IDENT, Lit: "Quantity"},
							Value: "Quantity",
						},
					},
				},
			},
			nil,
		},
		{
			`title
whereas the Amount in Stock (hereinafter Stock) is ninety-nine (99)
whereas the Current Quantity (hereinafter Quantity) is Stock
resolved publish Stock`,
			&ast.Resolution{
				WhereasStmts: []ast.WhereasStmt{
					&ast.DeclStmt{
						Token: token.Token{Typ: token.HEREINAFTER, Lit: "hereinafter"},
						Name: &ast.Identifier{
							Token: token.Token{Typ: token.IDENT, Lit: "Stock"},
							Value: "Stock",
						},
						Value: &ast.IntegerLiteral{
							Token: token.Token{Typ: token.INTEGER, Lit: "99"},
							Value: 99,
						},
					},
					&ast.DeclStmt{
						Token: token.Token{Typ: token.HEREINAFTER, Lit: "hereinafter"},
						Name: &ast.Identifier{
							Token: token.Token{Typ: token.IDENT, Lit: "Quantity"},
							Value: "Quantity",
						},
						Value: &ast.Identifier{
							Token: token.Token{Typ: token.IDENT, Lit: "Stock"},
							Value: "Stock",
						},
					},
				},
				ResolvedStmts: []ast.ResolvedStmt{
					&ast.PublishStmt{
						Token: token.Token{Typ: token.PUBLISH, Lit: "publish"},
						Value: &ast.Identifier{
							Token: token.Token{Typ: token.IDENT, Lit: "Stock"},
							Value: "Stock",
						},
					},
				},
			},
			unusedError{"Quantity"},
		},
	} {
		p := New(lexer.New(test.input))
		ast, err := p.ParseResolution()
		if err != nil {
			// Test the actual value of the last error generated
			err = p.errors[len(p.errors)-1]
		}
		if !reflect.DeepEqual(ast, test.ast) || err != test.err {
			t.Errorf("ParseResolution(%v): got %v, %v; want %v, %v", test.input, ast, err, test.ast, test.err)
		}
	}
}

func TestParseDeclStmt(t *testing.T) {
	for _, test := range []struct {
		input string
		want  *ast.DeclStmt
	}{
		{
			`hereinafter Greeting) is "Hello, World!"`,
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
		},
		{
			"hereinafter referred to as the Answer, is forty-two (42)",
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
		},
	} {
		p := New(lexer.New(test.input))
		got := p.parseDeclStmt()
		err := p.lastError()
		if err != nil || !reflect.DeepEqual(got, test.want) {
			t.Errorf("parseDeclStmt(%v): got %#v, %v, want %#v", test.input, got, err, test.want)
		}
	}
}

func TestParsePublishStmt(t *testing.T) {
	for _, test := range []struct {
		input string
		want  *ast.PublishStmt
	}{
		{
			`publish "Hello, World!"`,
			&ast.PublishStmt{
				Token: token.Token{Typ: token.PUBLISH, Lit: "publish"},
				Value: &ast.StringLiteral{
					Token: token.Token{Typ: token.STRING, Lit: "Hello, World!"},
					Value: "Hello, World!",
				},
			},
		},
		{
			"publish forty-two (42)",
			&ast.PublishStmt{
				Token: token.Token{Typ: token.PUBLISH, Lit: "publish"},
				Value: &ast.IntegerLiteral{
					Token: token.Token{Typ: token.INTEGER, Lit: "42"},
					Value: 42,
				},
			},
		},
		{
			"publish said Message",
			&ast.PublishStmt{
				Token: token.Token{Typ: token.PUBLISH, Lit: "publish"},
				Value: &ast.Identifier{
					Token: token.Token{Typ: token.IDENT, Lit: "Message"},
					Value: "Message",
				},
			},
		},
	} {
		p := New(lexer.New(test.input))
		p.idents["Message"] = declared
		got := p.parsePublishStmt()
		err := p.lastError()
		if err != nil || !reflect.DeepEqual(got, test.want) {
			t.Errorf("parsePublishStmt(%v): got %#v, %v, want %#v", test.input, got, err, test.want)
		}
	}
}

var identifierTests = []string{
	"Greeting",
	"Quantity",
	"Answer",
}

func TestParseIdentifier(t *testing.T) {
	for _, test := range identifierTests {
		want := &ast.Identifier{
			Token: token.Token{Typ: token.IDENT, Lit: test},
			Value: test,
		}
		p := New(lexer.New(test))
		if got := p.parseIdentifier(); !reflect.DeepEqual(got, want) {
			t.Errorf("parseIdentifier(%v): got %#v, want %#v", test, got, want)
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

func TestparseExpr(t *testing.T) {
	for _, test := range integerTests {
		input := fmt.Sprintf("%v (%v)", test.car, test.num)
		p := New(lexer.New(input))
		expr := p.parseExpr()
		err := p.lastError()
		if err != nil {
			t.Errorf("parseExpr(%v): got error %v", input, err)
		}
		e, ok := expr.(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("parseExpr(%v): got %T (%+v)", input, expr, expr)
		}
		if e.Value != test.n {
			t.Errorf("parseExpr(%v): got %v", input, e.Value)
		}
	}
	for _, test := range stringTests {
		input := fmt.Sprintf("\"%v\"", test)
		p := New(lexer.New(input))
		expr := p.parseExpr()
		err := p.lastError()
		if err != nil {
			t.Errorf("parseExpr(%v): got error %v", input, err)
		}
		e, ok := expr.(*ast.StringLiteral)
		if !ok {
			t.Errorf("parseExpr(%v): got %T (%+v)", input, expr, expr)
		}
		if e.Value != test {
			t.Errorf("parseExpr(%v): got %v", input, e.Value)
		}
	}
	for _, test := range identifierTests {
		p := New(lexer.New(test))
		expr := p.parseExpr()
		err := p.lastError()
		if err != nil {
			t.Errorf("parseExpr(%v): got error %v", test, err)
		}
		e, ok := expr.(*ast.Identifier)
		if !ok {
			t.Errorf("parseExpr(%v): got %T (%+v)", test, expr, expr)
		}
		if e.Value != test {
			t.Errorf("parseExpr(%v): got %v", test, e.Value)
		}
	}
}

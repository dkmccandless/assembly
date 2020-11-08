package parser

import (
	"fmt"
	"math"
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

func TestParseAssumeStmt(t *testing.T) {
	for _, test := range []struct {
		input string
		want  *ast.AssumeStmt
	}{
		{
			`Greeting assume "Greetings, Assembly."`,
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
		},
		{
			"Total assume three (3)",
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
		},
		{
			"Total assume sum Total one (1)",
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
		},
	} {
		p := New(lexer.New(test.input))
		p.idents["Greeting"] = declared
		p.idents["Total"] = declared
		got := p.parseResolvedStmt()
		err := p.lastError()
		if err != nil || !reflect.DeepEqual(got, test.want) {
			t.Errorf("parseAssumeStmt(%v): got %#v, %v, want %#v", test.input, got, err, test.want)
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

func TestParseExpr(t *testing.T) {
	for _, test := range []struct {
		input string
		expr  ast.Expr
	}{
		{"zero (0)", &ast.IntegerLiteral{token.Token{token.INTEGER, "0"}, 0}},
		{"one (1)", &ast.IntegerLiteral{token.Token{token.INTEGER, "1"}, 1}},
		{"negative three trillion (-3,000,000,000,000)", &ast.IntegerLiteral{token.Token{token.INTEGER, "-3000000000000"}, -3000000000000}},
		{
			"negative nine quintillion two hundred twenty-three quadrillion three hundred seventy-two trillion thirty-six billion eight hundred fifty-four million seven hundred seventy-five thousand eight hundred eight (-9,223,372,036,854,775,808)",
			&ast.IntegerLiteral{token.Token{token.INTEGER, "-9223372036854775808"}, math.MinInt64},
		},
		{
			"nine quintillion two hundred twenty-three quadrillion three hundred seventy-two trillion thirty-six billion eight hundred fifty-four million seven hundred seventy-five thousand eight hundred seven (9,223,372,036,854,775,807)",
			&ast.IntegerLiteral{token.Token{token.INTEGER, "9223372036854775807"}, math.MaxInt64},
		},

		{`""`, &ast.StringLiteral{token.Token{token.STRING, ""}, ""}},
		{`"WHEREAS"`, &ast.StringLiteral{token.Token{token.STRING, "WHEREAS"}, "WHEREAS"}},
		{`"zero (0)"`, &ast.StringLiteral{token.Token{token.STRING, "zero (0)"}, "zero (0)"}},
		{`"Greetings, Assembly."`, &ast.StringLiteral{token.Token{token.STRING, "Greetings, Assembly."}, "Greetings, Assembly."}},

		{"Greeting", &ast.Identifier{token.Token{token.IDENT, "Greeting"}, "Greeting"}},
		{"Quantity", &ast.Identifier{token.Token{token.IDENT, "Quantity"}, "Quantity"}},
		{"Answer", &ast.Identifier{token.Token{token.IDENT, "Answer"}, "Answer"}},

		// precedence tests
		{
			"ten (10) less thrice four (4)",
			&ast.InfixExpr{
				Token: token.Token{token.LESS, "less"},
				Left:  &ast.IntegerLiteral{token.Token{token.INTEGER, "10"}, 10},
				Right: &ast.UnaryPrefixExpr{
					Token: token.Token{token.THRICE, "thrice"},
					Right: &ast.IntegerLiteral{token.Token{token.INTEGER, "4"}, 4},
				},
			},
		},
		{
			"thrice four (4) less ten (10)",
			&ast.InfixExpr{
				Token: token.Token{token.LESS, "less"},
				Left: &ast.UnaryPrefixExpr{
					Token: token.Token{token.THRICE, "thrice"},
					Right: &ast.IntegerLiteral{token.Token{token.INTEGER, "4"}, 4},
				},
				Right: &ast.IntegerLiteral{token.Token{token.INTEGER, "10"}, 10},
			},
		},
		{
			"ten (10) less six (6) less one (1)",
			&ast.InfixExpr{
				Token: token.Token{token.LESS, "less"},
				Left: &ast.InfixExpr{
					Token: token.Token{token.LESS, "less"},
					Left:  &ast.IntegerLiteral{token.Token{token.INTEGER, "10"}, 10},
					Right: &ast.IntegerLiteral{token.Token{token.INTEGER, "6"}, 6},
				},
				Right: &ast.IntegerLiteral{token.Token{token.INTEGER, "1"}, 1},
			},
		},
		{
			"sum product Ax Bx product Ay By",
			&ast.BinaryPrefixExpr{
				Token: token.Token{token.SUM, "sum"},
				First: &ast.BinaryPrefixExpr{
					Token:  token.Token{token.PRODUCT, "product"},
					First:  &ast.Identifier{token.Token{token.IDENT, "Ax"}, "Ax"},
					Second: &ast.Identifier{token.Token{token.IDENT, "Bx"}, "Bx"},
				},
				Second: &ast.BinaryPrefixExpr{
					Token:  token.Token{token.PRODUCT, "product"},
					First:  &ast.Identifier{token.Token{token.IDENT, "Ay"}, "Ay"},
					Second: &ast.Identifier{token.Token{token.IDENT, "By"}, "By"},
				},
			},
		},
		{
			"product three (3) four (4) less two (2)",
			&ast.InfixExpr{
				Token: token.Token{token.LESS, "less"},
				Left: &ast.BinaryPrefixExpr{
					Token:  token.Token{token.PRODUCT, "product"},
					First:  &ast.IntegerLiteral{token.Token{token.INTEGER, "3"}, 3},
					Second: &ast.IntegerLiteral{token.Token{token.INTEGER, "4"}, 4},
				},
				Right: &ast.IntegerLiteral{token.Token{token.INTEGER, "2"}, 2},
			},
		},
		{
			"product three (3) less two (2) four (4)",
			&ast.BinaryPrefixExpr{
				Token: token.Token{token.PRODUCT, "product"},
				First: &ast.InfixExpr{
					Token: token.Token{token.LESS, "less"},
					Left:  &ast.IntegerLiteral{token.Token{token.INTEGER, "3"}, 3},
					Right: &ast.IntegerLiteral{token.Token{token.INTEGER, "2"}, 2},
				},
				Second: &ast.IntegerLiteral{token.Token{token.INTEGER, "4"}, 4},
			},
		},
		{
			"remainder twice eight (8) five (5)",
			&ast.BinaryPrefixExpr{
				Token: token.Token{token.REMAINDER, "remainder"},
				First: &ast.UnaryPrefixExpr{
					Token: token.Token{token.TWICE, "twice"},
					Right: &ast.IntegerLiteral{token.Token{token.INTEGER, "8"}, 8},
				},
				Second: &ast.IntegerLiteral{token.Token{token.INTEGER, "5"}, 5},
			},
		},
		{
			"twice remainder eight (8) five (5)",
			&ast.UnaryPrefixExpr{
				Token: token.Token{token.TWICE, "twice"},
				Right: &ast.BinaryPrefixExpr{
					Token:  token.Token{token.REMAINDER, "remainder"},
					First:  &ast.IntegerLiteral{token.Token{token.INTEGER, "8"}, 8},
					Second: &ast.IntegerLiteral{token.Token{token.INTEGER, "5"}, 5},
				},
			},
		},
		{
			"twice three (3) squared",
			&ast.UnaryPrefixExpr{
				Token: token.Token{token.TWICE, "twice"},
				Right: &ast.PostfixExpr{
					Token: token.Token{token.SQUARED, "squared"},
					Left:  &ast.IntegerLiteral{token.Token{token.INTEGER, "3"}, 3},
				},
			},
		},
		{
			"product two (2) ten (10) cubed",
			&ast.BinaryPrefixExpr{
				Token: token.Token{token.PRODUCT, "product"},
				First: &ast.IntegerLiteral{token.Token{token.INTEGER, "2"}, 2},
				Second: &ast.PostfixExpr{
					Token: token.Token{token.CUBED, "cubed"},
					Left:  &ast.IntegerLiteral{token.Token{token.INTEGER, "10"}, 10},
				},
			},
		},
		{
			"quotient seven (7) squared twelve (12)",
			&ast.BinaryPrefixExpr{
				Token: token.Token{token.QUOTIENT, "quotient"},
				First: &ast.PostfixExpr{
					Token: token.Token{token.SQUARED, "squared"},
					Left:  &ast.IntegerLiteral{token.Token{token.INTEGER, "7"}, 7},
				},
				Second: &ast.IntegerLiteral{token.Token{token.INTEGER, "12"}, 12},
			},
		},
		{
			"three (3) squared less two (2) cubed",
			&ast.InfixExpr{
				Token: token.Token{token.LESS, "less"},
				Left: &ast.PostfixExpr{
					Token: token.Token{token.SQUARED, "squared"},
					Left:  &ast.IntegerLiteral{token.Token{token.INTEGER, "3"}, 3},
				},
				Right: &ast.PostfixExpr{
					Token: token.Token{token.CUBED, "cubed"},
					Left:  &ast.IntegerLiteral{token.Token{token.INTEGER, "2"}, 2},
				},
			},
		},
	} {
		p := New(lexer.New(test.input))
		p.idents["Greeting"] = declared
		p.idents["Quantity"] = declared
		p.idents["Answer"] = declared
		p.idents["Ax"] = declared
		p.idents["Ay"] = declared
		p.idents["Bx"] = declared
		p.idents["By"] = declared
		expr := p.parseExpr(LOWEST)
		err := p.lastError()
		if !reflect.DeepEqual(expr, test.expr) || err != nil {
			t.Errorf("parseExpr(%v): got %+v, %v; want %+v", test.input, expr, err, test.expr)
		}
	}
}

func TestPostfixExpr(t *testing.T) {
	for _, test := range []struct {
		input string
		expr  ast.Expr
	}{
		{
			"three (3) squared",
			&ast.PostfixExpr{
				Token: token.Token{token.SQUARED, "squared"},
				Left:  &ast.IntegerLiteral{token.Token{token.INTEGER, "3"}, 3},
			},
		},
		{
			"four (4) cubed",
			&ast.PostfixExpr{
				Token: token.Token{token.CUBED, "cubed"},
				Left:  &ast.IntegerLiteral{token.Token{token.INTEGER, "4"}, 4},
			},
		},
		{
			"ten (10) cubed squared",
			&ast.PostfixExpr{
				Token: token.Token{token.SQUARED, "squared"},
				Left: &ast.PostfixExpr{
					Token: token.Token{token.CUBED, "cubed"},
					Left:  &ast.IntegerLiteral{token.Token{token.INTEGER, "10"}, 10},
				},
			},
		},
	} {
		p := New(lexer.New(test.input))
		expr := p.parseExpr(LOWEST)
		err := p.lastError()
		if !reflect.DeepEqual(expr, test.expr) || err != nil {
			t.Errorf("ParsePostfixExpr(%v): got %#v, %v; want %#v", test.input, expr, err, test.expr)
		}
	}
}

func TestParseUnaryPrefixExpr(t *testing.T) {
	for _, test := range []struct {
		input string
		expr  ast.Expr
	}{
		{
			"twice three (3)",
			&ast.UnaryPrefixExpr{
				Token: token.Token{token.TWICE, "twice"},
				Right: &ast.IntegerLiteral{token.Token{token.INTEGER, "3"}, 3},
			},
		},
		{
			"thrice four (4)",
			&ast.UnaryPrefixExpr{
				Token: token.Token{token.THRICE, "thrice"},
				Right: &ast.IntegerLiteral{token.Token{token.INTEGER, "4"}, 4},
			},
		},
		{
			"thrice twice negative one (-1)",
			&ast.UnaryPrefixExpr{
				Token: token.Token{token.THRICE, "thrice"},
				Right: &ast.UnaryPrefixExpr{
					Token: token.Token{token.TWICE, "twice"},
					Right: &ast.IntegerLiteral{token.Token{token.INTEGER, "-1"}, -1},
				},
			},
		},
	} {
		p := New(lexer.New(test.input))
		expr := p.parseExpr(LOWEST)
		err := p.lastError()
		if !reflect.DeepEqual(expr, test.expr) || err != nil {
			t.Errorf("ParseUnaryPrefixExpr(%v): got %#v, %v; want %#v", test.input, expr, err, test.expr)
		}
	}
}

func TestParseBinaryPrefixExpr(t *testing.T) {
	for _, test := range []struct {
		input string
		expr  ast.Expr
	}{
		{
			"sum one (1) one (1)",
			&ast.BinaryPrefixExpr{
				Token:  token.Token{token.SUM, "sum"},
				First:  &ast.IntegerLiteral{token.Token{token.INTEGER, "1"}, 1},
				Second: &ast.IntegerLiteral{token.Token{token.INTEGER, "1"}, 1},
			},
		},
		{
			"product two (2) three (3)",
			&ast.BinaryPrefixExpr{
				Token:  token.Token{token.PRODUCT, "product"},
				First:  &ast.IntegerLiteral{token.Token{token.INTEGER, "2"}, 2},
				Second: &ast.IntegerLiteral{token.Token{token.INTEGER, "3"}, 3},
			},
		},
		{
			"quotient twelve (12) five (5)",
			&ast.BinaryPrefixExpr{
				Token:  token.Token{token.QUOTIENT, "quotient"},
				First:  &ast.IntegerLiteral{token.Token{token.INTEGER, "12"}, 12},
				Second: &ast.IntegerLiteral{token.Token{token.INTEGER, "5"}, 5},
			},
		},
		{
			"remainder twelve (12) five (5)",
			&ast.BinaryPrefixExpr{
				Token:  token.Token{token.REMAINDER, "remainder"},
				First:  &ast.IntegerLiteral{token.Token{token.INTEGER, "12"}, 12},
				Second: &ast.IntegerLiteral{token.Token{token.INTEGER, "5"}, 5},
			},
		},
	} {
		p := New(lexer.New(test.input))
		expr := p.parseExpr(LOWEST)
		err := p.lastError()
		if !reflect.DeepEqual(expr, test.expr) || err != nil {
			t.Errorf("ParseBinaryPrefixExpr(%v): got %#v, %v; want %#v", test.input, expr, err, test.expr)
		}
	}
}

func TestParseInfixExpr(t *testing.T) {
	for _, test := range []struct {
		input string
		expr  ast.Expr
	}{
		{
			"three (3) less two (2)",
			&ast.InfixExpr{
				Token: token.Token{token.LESS, "less"},
				Left:  &ast.IntegerLiteral{token.Token{token.INTEGER, "3"}, 3},
				Right: &ast.IntegerLiteral{token.Token{token.INTEGER, "2"}, 2},
			},
		},
	} {
		p := New(lexer.New(test.input))
		expr := p.parseExpr(LOWEST)
		err := p.lastError()
		if !reflect.DeepEqual(expr, test.expr) || err != nil {
			t.Errorf("ParseInfixExpr(%v): got %+v, %v; want %+v", test.input, expr, err, test.expr)
		}
	}
}

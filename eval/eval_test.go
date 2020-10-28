package eval

import (
	"math"
	"testing"

	"github.com/dkmccandless/assembly/lexer"
	"github.com/dkmccandless/assembly/object"
	"github.com/dkmccandless/assembly/parser"
)

func TestEvalIntegerExpr(t *testing.T) {
	for _, test := range []struct {
		input string
		n     int64
	}{
		{"zero (0)", 0},
		{"one (1)", 1},
		{"ten (10)", 10},
		{"twenty (20)", 20},
		{"twenty-one (21)", 21},
		{"one hundred (100)", 100},
		{"one hundred one (101)", 101},
		{"one hundred ten (110)", 110},
		{"one hundred twenty (120)", 120},
		{"one hundred twenty-one (121)", 121},
		{"one thousand (1,000)", 1000},
		{"ten thousand (10,000)", 10000},
		{"twenty thousand (20,000)", 20000},
		{"twenty-one thousand (21,000)", 21000},
		{"one hundred thousand (100,000)", 100000},
		{"one hundred one thousand (101,000)", 101000},
		{"one hundred ten thousand (110,000)", 110000},
		{"one hundred twenty thousand (120,000)", 120000},
		{"one hundred twenty-one thousand (121,000)", 121000},
		{"one thousand one hundred twenty-one (1,121)", 1121},
		{"ten thousand one hundred twenty (10,120)", 10120},
		{"twenty thousand one hundred ten (20,110)", 20110},
		{"twenty-one thousand one hundred one (21,101)", 21101},
		{"one hundred thousand one hundred (100,100)", 100100},
		{"one hundred one thousand twenty-one (101,021)", 101021},
		{"one hundred ten thousand twenty (110,020)", 110020},
		{"one hundred twenty thousand ten (120,010)", 120010},
		{"one hundred twenty-one thousand one (121,001)", 121001},
		{"one million (1,000,000)", 1000000},
		{"one billion (1,000,000,000)", 1000000000},
		{"one trillion (1,000,000,000,000)", 1000000000000},
		{"one quadrillion (1,000,000,000,000,000)", 1000000000000000},
		{"one quintillion (1,000,000,000,000,000,000)", 1000000000000000000},
		{"one million one (1,000,001)", 1000001},
		{"one million one thousand (1,001,000)", 1001000},
		{"one billion one thousand (1,000,001,000)", 1000001000},
		{"one quintillion one (1,000,000,000,000,000,001)", 1000000000000000001},
		{"negative one (-1)", -1},
		{"negative ten (-10)", -10},
		{"negative twenty (-20)", -20},
		{"negative twenty-one (-21)", -21},
		{"negative one hundred (-100)", -100},
		{"negative one hundred one (-101)", -101},
		{"negative one hundred ten (-110)", -110},
		{"negative one hundred twenty (-120)", -120},
		{"negative one hundred twenty-one (-121)", -121},
		{"negative one thousand (-1,000)", -1000},
		{"negative ten thousand (-10,000)", -10000},
		{"negative twenty thousand (-20,000)", -20000},
		{"negative twenty-one thousand (-21,000)", -21000},
		{"negative one hundred thousand (-100,000)", -100000},
		{"negative one hundred one thousand (-101,000)", -101000},
		{"negative one hundred ten thousand (-110,000)", -110000},
		{"negative one hundred twenty thousand (-120,000)", -120000},
		{"negative one hundred twenty-one thousand (-121,000)", -121000},
		{"negative one thousand one hundred twenty-one (-1,121)", -1121},
		{"negative ten thousand one hundred twenty (-10,120)", -10120},
		{"negative twenty thousand one hundred ten (-20,110)", -20110},
		{"negative twenty-one thousand one hundred one (-21,101)", -21101},
		{"negative one hundred thousand one hundred (-100,100)", -100100},
		{"negative one hundred one thousand twenty-one (-101,021)", -101021},
		{"negative one hundred ten thousand twenty (-110,020)", -110020},
		{"negative one hundred twenty thousand ten (-120,010)", -120010},
		{"negative one hundred twenty-one thousand one (-121,001)", -121001},
		{"negative one million (-1,000,000)", -1000000},
		{"negative one billion (-1,000,000,000)", -1000000000},
		{"negative one trillion (-1,000,000,000,000)", -1000000000000},
		{"negative one quadrillion (-1,000,000,000,000,000)", -1000000000000000},
		{"negative one quintillion (-1,000,000,000,000,000,000)", -1000000000000000000},
		{"negative one million one (-1,000,001)", -1000001},
		{"negative one million one thousand (-1,001,000)", -1001000},
		{"negative one billion one thousand (-1,000,001,000)", -1000001000},
		{"negative one quintillion one (-1,000,000,000,000,000,001)", -1000000000000000001},
		{"negative nine quintillion two hundred twenty-three quadrillion three hundred seventy-two trillion thirty-six billion eight hundred fifty-four million seven hundred seventy-five thousand eight hundred eight (-9,223,372,036,854,775,808)", math.MinInt64},
		{"nine quintillion two hundred twenty-three quadrillion three hundred seventy-two trillion thirty-six billion eight hundred fifty-four million seven hundred seventy-five thousand eight hundred seven (9,223,372,036,854,775,807)", math.MaxInt64},
	} {
		node, err := parser.New(lexer.New(test.input)).ParseExpr()
		if err != nil {
			t.Errorf("EvalIntegerExpr(%v): got error %v", node, err)
		}
		obj, ok := Eval(node).(*object.Integer)
		if !ok {
			t.Errorf("EvalIntegerExpr(%v): got %T (%+v)", test.input, obj, obj)
		}
		if obj.Value != test.n {
			t.Errorf("EvalIntegerExpr(%v): got %v", test.input, obj.Value)
		}
	}
}

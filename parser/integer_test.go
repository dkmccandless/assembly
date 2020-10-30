package parser

import (
	"math"
	"reflect"
	"strconv"
	"testing"

	"github.com/dkmccandless/assembly/ast"
	"github.com/dkmccandless/assembly/lexer"
	"github.com/dkmccandless/assembly/token"
)

type integer struct {
	car, num string
	n        int64
}

// twoDigitIntegers contains all of the integers between 1 and 99 inclusive.
var twoDigitIntegers = []integer{
	{"one", "1", 1},
	{"two", "2", 2},
	{"three", "3", 3},
	{"four", "4", 4},
	{"five", "5", 5},
	{"six", "6", 6},
	{"seven", "7", 7},
	{"eight", "8", 8},
	{"nine", "9", 9},
	{"ten", "10", 10},
	{"eleven", "11", 11},
	{"twelve", "12", 12},
	{"thirteen", "13", 13},
	{"fourteen", "14", 14},
	{"fifteen", "15", 15},
	{"sixteen", "16", 16},
	{"seventeen", "17", 17},
	{"eighteen", "18", 18},
	{"nineteen", "19", 19},
	{"twenty", "20", 20},
	{"twenty-one", "21", 21},
	{"twenty-two", "22", 22},
	{"twenty-three", "23", 23},
	{"twenty-four", "24", 24},
	{"twenty-five", "25", 25},
	{"twenty-six", "26", 26},
	{"twenty-seven", "27", 27},
	{"twenty-eight", "28", 28},
	{"twenty-nine", "29", 29},
	{"thirty", "30", 30},
	{"thirty-one", "31", 31},
	{"thirty-two", "32", 32},
	{"thirty-three", "33", 33},
	{"thirty-four", "34", 34},
	{"thirty-five", "35", 35},
	{"thirty-six", "36", 36},
	{"thirty-seven", "37", 37},
	{"thirty-eight", "38", 38},
	{"thirty-nine", "39", 39},
	{"forty", "40", 40},
	{"forty-one", "41", 41},
	{"forty-two", "42", 42},
	{"forty-three", "43", 43},
	{"forty-four", "44", 44},
	{"forty-five", "45", 45},
	{"forty-six", "46", 46},
	{"forty-seven", "47", 47},
	{"forty-eight", "48", 48},
	{"forty-nine", "49", 49},
	{"fifty", "50", 50},
	{"fifty-one", "51", 51},
	{"fifty-two", "52", 52},
	{"fifty-three", "53", 53},
	{"fifty-four", "54", 54},
	{"fifty-five", "55", 55},
	{"fifty-six", "56", 56},
	{"fifty-seven", "57", 57},
	{"fifty-eight", "58", 58},
	{"fifty-nine", "59", 59},
	{"sixty", "60", 60},
	{"sixty-one", "61", 61},
	{"sixty-two", "62", 62},
	{"sixty-three", "63", 63},
	{"sixty-four", "64", 64},
	{"sixty-five", "65", 65},
	{"sixty-six", "66", 66},
	{"sixty-seven", "67", 67},
	{"sixty-eight", "68", 68},
	{"sixty-nine", "69", 69},
	{"seventy", "70", 70},
	{"seventy-one", "71", 71},
	{"seventy-two", "72", 72},
	{"seventy-three", "73", 73},
	{"seventy-four", "74", 74},
	{"seventy-five", "75", 75},
	{"seventy-six", "76", 76},
	{"seventy-seven", "77", 77},
	{"seventy-eight", "78", 78},
	{"seventy-nine", "79", 79},
	{"eighty", "80", 80},
	{"eighty-one", "81", 81},
	{"eighty-two", "82", 82},
	{"eighty-three", "83", 83},
	{"eighty-four", "84", 84},
	{"eighty-five", "85", 85},
	{"eighty-six", "86", 86},
	{"eighty-seven", "87", 87},
	{"eighty-eight", "88", 88},
	{"eighty-nine", "89", 89},
	{"ninety", "90", 90},
	{"ninety-one", "91", 91},
	{"ninety-two", "92", 92},
	{"ninety-three", "93", 93},
	{"ninety-four", "94", 94},
	{"ninety-five", "95", 95},
	{"ninety-six", "96", 96},
	{"ninety-seven", "97", 97},
	{"ninety-eight", "98", 98},
	{"ninety-nine", "99", 99},
}

// threeDigitIntegers contains all of the integers between 1 and 999 inclusive.
var threeDigitIntegers = func() []integer {
	td := append(make([]integer, 0, 999), twoDigitIntegers...)
	for _, hundred := range []integer{
		{"one hundred", "100", 100},
		{"two hundred", "200", 200},
		{"three hundred", "300", 300},
		{"four hundred", "400", 400},
		{"five hundred", "500", 500},
		{"six hundred", "600", 600},
		{"seven hundred", "700", 700},
		{"eight hundred", "800", 800},
		{"nine hundred", "900", 900},
	} {
		td = append(td, hundred)
		for _, twodigit := range twoDigitIntegers {
			td = append(td, integer{hundred.car + " " + twodigit.car, hundred.num[:3-len(twodigit.num)] + twodigit.num, hundred.n + twodigit.n})
		}
	}
	return td
}()

var integerTests = []integer{
	{"zero", "0", 0},

	{"one", "1", 1},                        // (nonzero) digit
	{"ten", "10", 10},                      // vigesimal
	{"twenty", "20", 20},                   // tens
	{"twenty-one", "21", 21},               // tens-digit
	{"one hundred", "100", 100},            // digit hundred
	{"one hundred one", "101", 101},        // digit hundred digit
	{"one hundred ten", "110", 110},        // digit hundred vigesimal
	{"one hundred twenty", "120", 120},     // digit hundred tens
	{"one hundred twenty-one", "121", 121}, // digit hundred tens-digit

	{"one thousand", "1,000", 1000},
	{"ten thousand", "10,000", 10000},
	{"twenty thousand", "20,000", 20000},
	{"twenty-one thousand", "21,000", 21000},
	{"one hundred thousand", "100,000", 100000},
	{"one hundred one thousand", "101,000", 101000},
	{"one hundred ten thousand", "110,000", 110000},
	{"one hundred twenty thousand", "120,000", 120000},
	{"one hundred twenty-one thousand", "121,000", 121000},

	{"one thousand one hundred twenty-one", "1,121", 1121},
	{"ten thousand one hundred twenty", "10,120", 10120},
	{"twenty thousand one hundred ten", "20,110", 20110},
	{"twenty-one thousand one hundred one", "21,101", 21101},
	{"one hundred thousand one hundred", "100,100", 100100},
	{"one hundred one thousand twenty-one", "101,021", 101021},
	{"one hundred ten thousand twenty", "110,020", 110020},
	{"one hundred twenty thousand ten", "120,010", 120010},
	{"one hundred twenty-one thousand one", "121,001", 121001},

	{"one million", "1,000,000", 1000000},
	{"one billion", "1,000,000,000", 1000000000},
	{"one trillion", "1,000,000,000,000", 1000000000000},
	{"one quadrillion", "1,000,000,000,000,000", 1000000000000000},
	{"one quintillion", "1,000,000,000,000,000,000", 1000000000000000000},

	{"one million one", "1,000,001", 1000001},
	{"one million one thousand", "1,001,000", 1001000},
	{"one billion one thousand", "1,000,001,000", 1000001000},
	{"one quintillion one", "1,000,000,000,000,000,001", 1000000000000000001},

	{"negative one", "-1", -1},
	{"negative ten", "-10", -10},
	{"negative twenty", "-20", -20},
	{"negative twenty-one", "-21", -21},
	{"negative one hundred", "-100", -100},
	{"negative one hundred one", "-101", -101},
	{"negative one hundred ten", "-110", -110},
	{"negative one hundred twenty", "-120", -120},
	{"negative one hundred twenty-one", "-121", -121},

	{"negative one thousand", "-1,000", -1000},
	{"negative ten thousand", "-10,000", -10000},
	{"negative twenty thousand", "-20,000", -20000},
	{"negative twenty-one thousand", "-21,000", -21000},
	{"negative one hundred thousand", "-100,000", -100000},
	{"negative one hundred one thousand", "-101,000", -101000},
	{"negative one hundred ten thousand", "-110,000", -110000},
	{"negative one hundred twenty thousand", "-120,000", -120000},
	{"negative one hundred twenty-one thousand", "-121,000", -121000},

	{"negative one thousand one hundred twenty-one", "-1,121", -1121},
	{"negative ten thousand one hundred twenty", "-10,120", -10120},
	{"negative twenty thousand one hundred ten", "-20,110", -20110},
	{"negative twenty-one thousand one hundred one", "-21,101", -21101},
	{"negative one hundred thousand one hundred", "-100,100", -100100},
	{"negative one hundred one thousand twenty-one", "-101,021", -101021},
	{"negative one hundred ten thousand twenty", "-110,020", -110020},
	{"negative one hundred twenty thousand ten", "-120,010", -120010},
	{"negative one hundred twenty-one thousand one", "-121,001", -121001},

	{"negative one million", "-1,000,000", -1000000},
	{"negative one billion", "-1,000,000,000", -1000000000},
	{"negative one trillion", "-1,000,000,000,000", -1000000000000},
	{"negative one quadrillion", "-1,000,000,000,000,000", -1000000000000000},
	{"negative one quintillion", "-1,000,000,000,000,000,000", -1000000000000000000},

	{"negative one million one", "-1,000,001", -1000001},
	{"negative one million one thousand", "-1,001,000", -1001000},
	{"negative one billion one thousand", "-1,000,001,000", -1000001000},
	{"negative one quintillion one", "-1,000,000,000,000,000,001", -1000000000000000001},

	{
		"negative nine quintillion two hundred twenty-three quadrillion three hundred seventy-two trillion thirty-six billion eight hundred fifty-four million seven hundred seventy-five thousand eight hundred eight",
		"-9,223,372,036,854,775,808",
		math.MinInt64,
	},
	{
		"nine quintillion two hundred twenty-three quadrillion three hundred seventy-two trillion thirty-six billion eight hundred fifty-four million seven hundred seventy-five thousand eight hundred seven",
		"9,223,372,036,854,775,807",
		math.MaxInt64,
	},
}

func TestParseTwoDigitCardinal(t *testing.T) {
	for _, test := range twoDigitIntegers {
		p := New(lexer.New(test.car))
		if got, err := p.parseTwoDigitCardinal(); got != test.n || err != nil {
			t.Errorf("parseTwoDigitCardinal(%v): got %v, %v; want %v", test.car, got, err, test.n)
		}
	}
}

func TestParseThreeDigitCardinal(t *testing.T) {
	for _, test := range threeDigitIntegers {
		p := New(lexer.New(test.car))
		if got, err := p.parseThreeDigitCardinal(); got != test.n || err != nil {
			t.Errorf("parseThreeDigitCardinal(%v): got %v, %v; want %v", test.car, got, err, test.n)
		}
	}
}

func TestParseCardinalLiteral(t *testing.T) {
	for _, test := range integerTests {
		p := New(lexer.New(test.car))
		if got, err := p.parseCardinalLiteral(); got != test.n || err != nil {
			t.Errorf("parseCardinalLiteral(%v): got %v, %v; want %v", test.car, got, err, test.n)
		}
	}
}

func TestParseInvalidCardinalLiteral(t *testing.T) {
	for _, car := range []string{
		"",
		"0",
		"hundred",
		"thousand",
		"twenty-",
		"negative 1",
		"sixty-twelve",
		"twenty-thousand",
		"one thousand one million",
		"one million one million",
		"one billion one thousand one million",

		// overflow
		"negative nine quintillion two hundred twenty-three quadrillion three hundred seventy-two trillion thirty-six billion eight hundred fifty-four million seven hundred seventy-five thousand eight hundred nine",
		"nine quintillion two hundred twenty-three quadrillion three hundred seventy-two trillion thirty-six billion eight hundred fifty-four million seven hundred seventy-five thousand eight hundred eight",
	} {
		p := New(lexer.New(car))
		if got, err := p.parseCardinalLiteral(); got != 0 || err != errCardinal {
			t.Errorf("parseCardinalLiteral(%v): got %v, %v; want 0, %v", car, got, err, errCardinal)
		}
	}
}

func TestParseNumeralLiteral(t *testing.T) {
	for _, test := range integerTests {
		p := New(lexer.New(test.num))
		if got, err := p.parseNumeralLiteral(); got != test.n || err != nil {
			t.Errorf("parseNumeralLiteral(%v): got %v, %v; want %v", test.num, got, err, test.n)
		}
	}
}

func TestParseInvalidNumeralLiteral(t *testing.T) {
	for _, num := range []string{
		"",
		"-",
		"a",
		"-0", // zero (0) is neither positive nor negative
		"00",
		"01",
		"1000",
		",100",
		"10,00",
		"-,100",
		"1,000,",
		"1,,000",
		"1000,000",

		// overflow
		"-9,223,372,036,854,775,809",
		"9,223,372,036,854,775,808",
	} {
		p := New(lexer.New(num))
		if got, err := p.parseNumeralLiteral(); got != 0 || err != errNumeral {
			t.Errorf("parseNumeralLiteral(%v): got %v, %v; want 0, %v", num, got, err, errNumeral)
		}
	}
}

func TestParseIntegerLiteral(t *testing.T) {
	for _, test := range integerTests {
		input := test.car + " (" + test.num + ")"
		want := &ast.IntegerLiteral{
			Token: token.Token{
				Typ: token.INTEGER,
				Lit: strconv.Itoa(int(test.n)),
			},
			Value: test.n,
		}
		p := New(lexer.New(input))
		if got, err := p.parseIntegerLiteral(); !reflect.DeepEqual(got, want) || err != nil {
			t.Errorf("parseIntegerLiteral(%v): got %v, %v; want %v", input, got, err, test.n)
		}
	}
}

func TestParseInvalidIntegerLiteral(t *testing.T) {
	for _, test := range []struct {
		input string
		want  error
	}{
		{"0", errInteger},                               // numeral without cardinal
		{"zero", errInteger},                            // cardinal without numeral
		{"one 1", errInteger},                           // numeral must be parenthesized
		{"two 2)", errInteger},                          // parentheses must be properly opened...
		{"two (2", errInteger},                          // ...and closed
		{"three ($3)", errInteger},                      // parentheses may only contain a numeral
		{"negative one (1)", errDisagree},               // parity of cardinal and numeral must agree
		{"three hundred sixty-five (366)", errDisagree}, // magnitude of cardinal and numeral must agree
	} {
		p := New(lexer.New(test.input))
		if got, err := p.parseIntegerLiteral(); got != nil || err != test.want {
			t.Errorf("parseIntegerLiteral(%v): got %v, %v; want %v", test.input, got, err, test.want)
		}
	}
}

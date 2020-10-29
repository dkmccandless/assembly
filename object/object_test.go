package object

import (
	"math"
	"testing"
)

var integerTests = []struct {
	n int64
	s string
}{
	{0, "zero (0)"},
	{1, "one (1)"},
	{10, "ten (10)"},
	{20, "twenty (20)"},
	{21, "twenty-one (21)"},
	{100, "one hundred (100)"},
	{101, "one hundred one (101)"},
	{110, "one hundred ten (110)"},
	{120, "one hundred twenty (120)"},
	{121, "one hundred twenty-one (121)"},
	{1000, "one thousand (1,000)"},
	{10000, "ten thousand (10,000)"},
	{20000, "twenty thousand (20,000)"},
	{21000, "twenty-one thousand (21,000)"},
	{100000, "one hundred thousand (100,000)"},
	{101000, "one hundred one thousand (101,000)"},
	{110000, "one hundred ten thousand (110,000)"},
	{120000, "one hundred twenty thousand (120,000)"},
	{121000, "one hundred twenty-one thousand (121,000)"},
	{1121, "one thousand one hundred twenty-one (1,121)"},
	{10120, "ten thousand one hundred twenty (10,120)"},
	{20110, "twenty thousand one hundred ten (20,110)"},
	{21101, "twenty-one thousand one hundred one (21,101)"},
	{100100, "one hundred thousand one hundred (100,100)"},
	{101021, "one hundred one thousand twenty-one (101,021)"},
	{110020, "one hundred ten thousand twenty (110,020)"},
	{120010, "one hundred twenty thousand ten (120,010)"},
	{121001, "one hundred twenty-one thousand one (121,001)"},
	{1000000, "one million (1,000,000)"},
	{1000000000, "one billion (1,000,000,000)"},
	{1000000000000, "one trillion (1,000,000,000,000)"},
	{1000000000000000, "one quadrillion (1,000,000,000,000,000)"},
	{1000000000000000000, "one quintillion (1,000,000,000,000,000,000)"},
	{1000001, "one million one (1,000,001)"},
	{1001000, "one million one thousand (1,001,000)"},
	{1000001000, "one billion one thousand (1,000,001,000)"},
	{1000000000000000001, "one quintillion one (1,000,000,000,000,000,001)"},
	{-1, "negative one (-1)"},
	{-10, "negative ten (-10)"},
	{-20, "negative twenty (-20)"},
	{-21, "negative twenty-one (-21)"},
	{-100, "negative one hundred (-100)"},
	{-101, "negative one hundred one (-101)"},
	{-110, "negative one hundred ten (-110)"},
	{-120, "negative one hundred twenty (-120)"},
	{-121, "negative one hundred twenty-one (-121)"},
	{-1000, "negative one thousand (-1,000)"},
	{-10000, "negative ten thousand (-10,000)"},
	{-20000, "negative twenty thousand (-20,000)"},
	{-21000, "negative twenty-one thousand (-21,000)"},
	{-100000, "negative one hundred thousand (-100,000)"},
	{-101000, "negative one hundred one thousand (-101,000)"},
	{-110000, "negative one hundred ten thousand (-110,000)"},
	{-120000, "negative one hundred twenty thousand (-120,000)"},
	{-121000, "negative one hundred twenty-one thousand (-121,000)"},
	{-1121, "negative one thousand one hundred twenty-one (-1,121)"},
	{-10120, "negative ten thousand one hundred twenty (-10,120)"},
	{-20110, "negative twenty thousand one hundred ten (-20,110)"},
	{-21101, "negative twenty-one thousand one hundred one (-21,101)"},
	{-100100, "negative one hundred thousand one hundred (-100,100)"},
	{-101021, "negative one hundred one thousand twenty-one (-101,021)"},
	{-110020, "negative one hundred ten thousand twenty (-110,020)"},
	{-120010, "negative one hundred twenty thousand ten (-120,010)"},
	{-121001, "negative one hundred twenty-one thousand one (-121,001)"},
	{-1000000, "negative one million (-1,000,000)"},
	{-1000000000, "negative one billion (-1,000,000,000)"},
	{-1000000000000, "negative one trillion (-1,000,000,000,000)"},
	{-1000000000000000, "negative one quadrillion (-1,000,000,000,000,000)"},
	{-1000000000000000000, "negative one quintillion (-1,000,000,000,000,000,000)"},
	{-1000001, "negative one million one (-1,000,001)"},
	{-1001000, "negative one million one thousand (-1,001,000)"},
	{-1000001000, "negative one billion one thousand (-1,000,001,000)"},
	{-1000000000000000001, "negative one quintillion one (-1,000,000,000,000,000,001)"},
	{math.MinInt64, "negative nine quintillion two hundred twenty-three quadrillion three hundred seventy-two trillion thirty-six billion eight hundred fifty-four million seven hundred seventy-five thousand eight hundred eight (-9,223,372,036,854,775,808)"},
	{math.MaxInt64, "nine quintillion two hundred twenty-three quadrillion three hundred seventy-two trillion thirty-six billion eight hundred fifty-four million seven hundred seventy-five thousand eight hundred seven (9,223,372,036,854,775,807)"},
}

func TestIntegerInspect(t *testing.T) {
	for _, test := range integerTests {
		i := &Integer{Value: test.n}
		if got := i.Inspect(); got != test.s {
			t.Errorf("Inspect(%v): got %v, want %v", test.n, got, test.s)
		}
	}
}

var stringTests = []string{
	"",
	"WHEREAS",
	"zero (0)",
	"Greetings, Assembly.",
}

func TestStringInspect(t *testing.T) {
	for _, test := range stringTests {
		s := &String{Value: test}
		if got := s.Inspect(); got != test {
			t.Errorf("Inspect(%v): got %v", test, got)
		}
	}
}

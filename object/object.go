package object

import "fmt"

type Object interface {
	Type() Type
	Inspect() string
}

type Type int

const (
	INTEGER Type = iota
	STRING
	ERROR
)

type Integer struct{ Value int64 }

func (i *Integer) Type() Type { return INTEGER }
func (i *Integer) Inspect() string {
	n := i.Value
	if n == 0 {
		return "zero (0)"
	}
	var car, num string
	negative := n < 0

	var groups []int64
	for n != 0 {
		// Negate n piecewise to avoid overflow in the case of math.MinInt64, for which n == -n
		r := n % 1000
		if negative {
			r *= -1
		}
		groups = append(groups, r)
		n /= 1000
	}
	for i := len(groups) - 1; i >= 0; i-- {
		n := groups[i]
		if i == len(groups)-1 {
			num += fmt.Sprintf("%d", n)
		} else {
			num += fmt.Sprintf(",%03d", n)
		}
		if n == 0 {
			continue
		}
		if len(car) > 0 {
			car += " "
		}
		var (
			ones      = []string{"", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
			vigesimal = []string{"ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen"}
			tens      = []string{"", "", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety"}
			powers    = []string{"", "thousand", "million", "billion", "trillion", "quadrillion", "quintillion"}
		)
		if n >= 100 {
			car += ones[n/100] + " hundred"
			n %= 100
			if n > 0 {
				car += " "
			}
		}
		switch {
		case n == 0:
		case n < 10:
			car += ones[n]
		case n < 20:
			car += vigesimal[n-10]
		default:
			car += tens[n/10]
			if n%10 != 0 {
				car += "-" + ones[n%10]
			}
		}
		if i > 0 {
			car += " " + powers[i]
		}
	}

	if negative {
		car, num = "negative "+car, "-"+num
	}

	return fmt.Sprintf("%v (%v)", car, num)
}

type String struct{ Value string }

func (s *String) Type() Type      { return STRING }
func (s *String) Inspect() string { return s.Value }

type Error struct{ Value string }

func (e *Error) Type() Type      { return ERROR }
func (e *Error) Inspect() string { return e.Value }

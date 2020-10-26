package parser

import (
	"errors"
	"strconv"

	"github.com/dkmccandless/assembly/ast"
	"github.com/dkmccandless/assembly/token"
)

var (
	// errInteger is returned by parseInteger given input that does not consist of a cardinal followed by a parenthesized numeral.
	errInteger = errors.New("invalid integer")

	// errCardinal is returned by cardinal parsing functions given input that cannot be parsed.
	errCardinal = errors.New("invalid cardinal")

	// errNumeral is returned by numeral parsing functions given input that cannot be parsed.
	errNumeral = errors.New("invalid numeral")

	// errDisagree is returned by parseIntegerLit when a cardinal and the following numeral represent different numbers.
	errDisagree = errors.New("cardinal and numeral disagree")
)

func (p *Parser) parseIntegerLiteral() (ast.Expr, error) {
	c, err := p.parseCardinalLiteral()
	if err != nil {
		return nil, err
	}

	if !p.peekIs(token.LPAREN) {
		return nil, errInteger
	}
	p.next()

	if !p.peekIs(token.NUMERAL) && !p.peekIs(token.DASH) {
		return nil, errInteger
	}
	p.next()

	n, err := p.parseNumeralLiteral()
	if err != nil {
		return nil, err
	}

	if !p.peekIs(token.RPAREN) {
		return nil, errInteger
	}
	p.next()

	if c != n {
		return nil, errDisagree
	}

	return &ast.IntegerLiteral{Token: token.Token{token.INTEGER, strconv.Itoa(int(n))}, Value: n}, nil
}

func (p *Parser) parseCardinalLiteral() (int64, error) {
	if p.curIs(token.ZERO) {
		return 0, nil
	}

	var n int64
	var negative bool
	if p.curIs(token.NEGATIVE) {
		negative = true
		p.next()
	}

	for {
		td, err := p.parseThreeDigitCardinal()
		if err != nil {
			return 0, err
		}
		if negative {
			td *= -1
		}

		var pow int64 = 1
		if p.peekIs(token.POWER) {
			p.next()
			pow = value[p.cur.Lit]
		}

		// Check for direct overflow
		if pow == 1e18 && (td > 9 || td < -9) {
			return 0, errCardinal
		}

		// Check that pow is smaller than all preceding powers
		if pow < 1e18 && n%(pow*1000) != 0 {
			return 0, errCardinal
		}

		// Add and check for overflow
		oldn := n
		n += td * pow
		if n == oldn || (n > oldn) == negative {
			return 0, errCardinal
		}

		if pow == 1 || !p.peekIs(token.ONES) && !p.peekIs(token.VIGESIMAL) && !p.peekIs(token.TENS) {
			return n, nil
		}
		p.next()
	}
}

func (p *Parser) parseThreeDigitCardinal() (int64, error) {
	switch p.cur.Typ {
	case token.ONES:
		n := value[p.cur.Lit]
		if p.peek.Typ == token.HUNDRED {
			p.next()
			n *= 100
			if p.peekIs(token.ONES) || p.peekIs(token.VIGESIMAL) || p.peekIs(token.TENS) {
				p.next()
				td, err := p.parseTwoDigitCardinal()
				if err != nil {
					return 0, err
				}
				n += td
			}
		}
		return n, nil
	case token.VIGESIMAL, token.TENS:
		return p.parseTwoDigitCardinal()
	default:
		return 0, errCardinal
	}
}

func (p *Parser) parseTwoDigitCardinal() (int64, error) {
	switch val := value[p.cur.Lit]; p.cur.Typ {
	case token.ONES, token.VIGESIMAL:
		return val, nil
	case token.TENS:
		n := val
		if p.peekIs(token.DASH) {
			p.next()
			if !p.peekIs(token.ONES) {
				return 0, errCardinal
			}
			p.next()
			n += value[p.cur.Lit]
		}
		return n, nil
	default:
		return 0, errCardinal
	}
}

var value = map[string]int64{
	"one":         1,
	"two":         2,
	"three":       3,
	"four":        4,
	"five":        5,
	"six":         6,
	"seven":       7,
	"eight":       8,
	"nine":        9,
	"ten":         10,
	"eleven":      11,
	"twelve":      12,
	"thirteen":    13,
	"fourteen":    14,
	"fifteen":     15,
	"sixteen":     16,
	"seventeen":   17,
	"eighteen":    18,
	"nineteen":    19,
	"twenty":      20,
	"thirty":      30,
	"forty":       40,
	"fifty":       50,
	"sixty":       60,
	"seventy":     70,
	"eighty":      80,
	"ninety":      90,
	"hundred":     100,
	"thousand":    1e3,
	"million":     1e6,
	"billion":     1e9,
	"trillion":    1e12,
	"quadrillion": 1e15,
	"quintillion": 1e18,
}

func (p *Parser) parseNumeralLiteral() (int64, error) {
	num := p.cur.Lit
	if len(num) == 0 {
		return 0, errNumeral
	}
	if num == "0" {
		return 0, nil
	}

	var n int64
	var negative bool
	if num[0] == '-' {
		negative = true
		num = num[1:]
		if len(num) == 0 {
			return 0, errNumeral
		}
	}

	if num[0] < '1' || '9' < num[0] {
		return 0, errNumeral
	}

	for i, b := range num {
		// Commas must delimit every group of three (3) digits.
		if i%4 == len(num)%4 {
			if b != ',' {
				return 0, errNumeral
			}
			continue
		}
		if b < '0' || '9' < b {
			return 0, errNumeral
		}

		d := int64(b - '0')
		if negative {
			d *= -1
		}

		// Add and check for overflow
		oldn := n
		n = 10*n + d
		if n == oldn || (n > oldn) == negative {
			return 0, errNumeral
		}
	}

	return n, nil
}

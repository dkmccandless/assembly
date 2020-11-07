package token

import "strings"

// Token is a lexical token of Assembly source code.
type Token struct {
	Typ Type
	Lit string
}

// Type is a token type.
type Type int

const (
	// Special tokens
	EOF Type = iota
	COMMENT

	// Identifiers and type literals
	IDENT
	INTEGER
	NUMERAL
	STRING

	// Cardinals
	NEGATIVE
	ZERO
	ONES
	VIGESIMAL
	TENS
	HUNDRED
	POWER

	// Operators
	TWICE
	THRICE
	LESS

	// Punctuation
	LPAREN
	RPAREN
	DASH

	// Keywords
	WHEREAS
	RESOLVED
	HEREINAFTER
	PUBLISH
)

var keywords = map[string]Type{
	"negative":    NEGATIVE,
	"zero":        ZERO,
	"one":         ONES,
	"two":         ONES,
	"three":       ONES,
	"four":        ONES,
	"five":        ONES,
	"six":         ONES,
	"seven":       ONES,
	"eight":       ONES,
	"nine":        ONES,
	"ten":         VIGESIMAL,
	"eleven":      VIGESIMAL,
	"twelve":      VIGESIMAL,
	"thirteen":    VIGESIMAL,
	"fourteen":    VIGESIMAL,
	"fifteen":     VIGESIMAL,
	"sixteen":     VIGESIMAL,
	"seventeen":   VIGESIMAL,
	"eighteen":    VIGESIMAL,
	"nineteen":    VIGESIMAL,
	"twenty":      TENS,
	"thirty":      TENS,
	"forty":       TENS,
	"fifty":       TENS,
	"sixty":       TENS,
	"seventy":     TENS,
	"eighty":      TENS,
	"ninety":      TENS,
	"hundred":     HUNDRED,
	"thousand":    POWER,
	"million":     POWER,
	"billion":     POWER,
	"trillion":    POWER,
	"quadrillion": POWER,
	"quintillion": POWER,

	"twice":  TWICE,
	"thrice": THRICE,
	"less":   LESS,

	"whereas":     WHEREAS,
	"resolved":    RESOLVED,
	"hereinafter": HEREINAFTER,
	"publish":     PUBLISH,
}

// Lookup maps s to its keyword Type, if any,
// or else to IDENT if it begins with a capital letter
// or COMMENT otherwise.
func Lookup(s string) Type {
	lower := strings.ToLower(s)
	if typ, ok := keywords[lower]; ok {
		return typ
	}
	if 'A' <= s[0] && s[0] <= 'Z' {
		return IDENT
	}
	return COMMENT
}

// IsCardinal reports whether t's Type is one of the cardinal numeric Types.
func (t Token) IsCardinal() bool {
	return t.Typ == NEGATIVE ||
		t.Typ == ZERO ||
		t.Typ == ONES ||
		t.Typ == VIGESIMAL ||
		t.Typ == TENS ||
		t.Typ == HUNDRED ||
		t.Typ == POWER
}

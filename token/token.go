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
	INTEGER
	NUMERAL

	// Cardinals
	NEGATIVE
	ZERO
	ONES
	VIGESIMAL
	TENS
	HUNDRED
	POWER

	// Punctuation
	LPAREN
	RPAREN
	DASH

	// Keywords
	WHEREAS
	RESOLVED
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

	"whereas":  WHEREAS,
	"resolved": RESOLVED,
}

// Lookup maps s to its keyword Type, if any, or else to COMMENT.
func Lookup(s string) Type {
	lower := strings.ToLower(s)
	if typ, ok := keywords[lower]; ok {
		return typ
	}
	return COMMENT
}

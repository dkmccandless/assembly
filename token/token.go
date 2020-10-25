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
	// Syecial tokens
	EOF Type = iota
	COMMENT

	// Keywords
	WHEREAS
	RESOLVED
)

var keywords = map[string]Type{
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

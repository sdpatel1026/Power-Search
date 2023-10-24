package tokenize

import (
	"strings"
	"unicode"
)

type EnTokenizer struct {
	Seprators string
}

// Tokens splits text into tokens.
func (enTokenizer *EnTokenizer) Tokens(text string) []string {
	splitter := func(r rune) bool {
		if unicode.IsSpace(r) {
			return true
		}
		return strings.ContainsRune(enTokenizer.Seprators, r)
	}
	return strings.FieldsFunc(text, splitter)
}

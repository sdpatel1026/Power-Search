package tfidf

import (
	"strings"
)

func cleanWord(word string) string {

	word = strings.TrimSpace(word)
	word = strings.Trim(word, ":-")
	word = strings.Trim(word, "-")
	word = strings.Trim(word, "—")
	word = strings.Trim(word, ":")
	word = strings.Trim(word, ".")
	word = strings.Trim(word, ",")
	word = strings.Trim(word, "‘")
	word = strings.Trim(word, "’")
	word = strings.Trim(word, `"`)
	word = strings.Trim(word, `“`)
	word = strings.Trim(word, `”`)
	word = strings.Trim(word, `**`)
	word = strings.Trim(word, `##`)
	word = strings.Trim(word, `*`)
	word = strings.Trim(word, `#`)
	word = strings.Trim(word, "'")
	word = strings.Trim(word, ")")
	word = strings.Trim(word, "(")
	word = strings.Trim(word, "[")
	word = strings.Trim(word, "]")
	word = strings.Trim(word, "{")
	word = strings.Trim(word, "}")
	return word
}

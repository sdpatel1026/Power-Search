package tokenize

type Tokenizer interface {
	Tokens(text string) []string
}

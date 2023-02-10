package tokenizer

import (
	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/model/wordlevel"
)

type RuneLevelVocab interface {
	Vocab() map[string]int
	UnkToken() string
	TokenToID(token string) (int, error)
}

// RuneLevel is a model tokenizer that splits each word into runes and maps runes to IDs.
type RuneLevel struct {
	*wordlevel.WordLevel

	vocab RuneLevelVocab
}

func NewRuneLevel(vocab RuneLevelVocab) *RuneLevel {
	return &RuneLevel{
		WordLevel: NewWordLevel(vocab.Vocab(), vocab.UnkToken()),
		vocab:     vocab,
	}
}

// Tokenize transforms given input token into a list of rune-level sub-tokens.
func (rl *RuneLevel) Tokenize(token string) ([]tokenizer.Token, error) {
	var tokens []tokenizer.Token

	var offset int
	for _, r := range []rune(token) {
		s := string(r)

		id, err := rl.vocab.TokenToID(s)
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, tokenizer.Token{
			Id:      id,
			Value:   s,
			Offsets: []int{offset, offset + len(s)},
		})

		offset += len(s)
	}
	return tokens, nil
}

// NewWordLevel creates a WordLevel model from a given vocab.
func NewWordLevel(vocab map[string]int, unkToken string) *wordlevel.WordLevel {
	builder := wordlevel.NewWordLevelBuilder()
	builder.Vocab(vocab)
	builder.UnkToken(unkToken)
	return builder.Build()
}

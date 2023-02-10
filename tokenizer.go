package tokenizer

import (
	"github.com/sugarme/tokenizer"
)

type Tokenizer struct {
	*tokenizer.Tokenizer
}

func (t *Tokenizer) EncodeBatchTexts(texts []string, addSpecialTokens bool) ([]tokenizer.Encoding, error) {
	var inputs []tokenizer.EncodeInput
	for _, text := range texts {
		seq := tokenizer.NewInputSequence(text)
		inputs = append(inputs, tokenizer.NewSingleEncodeInput(seq))
	}
	return t.EncodeBatchSerially(inputs, addSpecialTokens)
}

// EncodeBatchSerially encodes all sentences serially.
func (t *Tokenizer) EncodeBatchSerially(inputs []tokenizer.EncodeInput, addSpecialTokens bool) ([]tokenizer.Encoding, error) {
	var encodings []tokenizer.Encoding
	for _, input := range inputs {
		e, err := t.Tokenizer.Encode(input, addSpecialTokens)
		if err != nil {
			return nil, err
		}
		encodings = append(encodings, *e)
	}

	// Do padding if specified.
	padding := t.Tokenizer.GetPadding()
	if padding != nil {
		encodings = tokenizer.PadEncodings(encodings, *padding)
	}

	return encodings, nil
}

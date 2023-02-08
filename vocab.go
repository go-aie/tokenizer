package tokenizer

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

type Vocab[T constraints.Integer] struct {
	Vocab    map[string]T
	InvVocab map[T]string
}

func NewVocabFromFile[T constraints.Integer](filename string, separator string) (*Vocab[T], error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return NewVocabFromSlice[T](lines, separator)
}

func NewVocabFromSlice[T constraints.Integer](lines []string, separator string) (*Vocab[T], error) {
	var vocab = make(map[string]T)
	var invVocab = make(map[T]string)

	for i, line := range lines {
		text := strings.TrimSpace(line)
		sub := strings.Split(text, separator)

		var token string
		var id T

		switch len(sub) {
		case 1:
			// sub: <token>
			token, id = sub[0], T(i)
		case 2:
			if n, err := strconv.Atoi(sub[0]); err == nil {
				// sub: <id> <token>
				id, token = T(n), sub[1]
			} else {
				n, err = strconv.Atoi(sub[1])
				if err != nil {
					return nil, fmt.Errorf("invalid content: %q at line %d", line, i+1)
				}

				// sub: <token> <id>
				token, id = sub[0], T(n)
			}
		default:
			return nil, fmt.Errorf("invalid content: %q at line %d", line, i+1)
		}

		vocab[token] = id
		invVocab[id] = token
	}

	return &Vocab[T]{
		Vocab:    vocab,
		InvVocab: invVocab,
	}, nil
}

func (v *Vocab[T]) TokensToIDs(tokens []string) (ids []T) {
	for _, token := range tokens {
		if id, ok := v.Vocab[token]; ok {
			ids = append(ids, id)
		}
	}
	return
}

func (v *Vocab[T]) IDsToTokens(ids []T) (tokens []string) {
	for _, id := range ids {
		if token, ok := v.InvVocab[id]; ok {
			tokens = append(tokens, token)
		}
	}
	return
}

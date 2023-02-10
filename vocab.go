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
	vocab    map[string]T
	invVocab map[T]string
	unkToken string
}

func NewVocabFromFile[T constraints.Integer](filename, separator, unkToken string) (*Vocab[T], error) {
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

	return NewVocabFromSlice[T](lines, separator, unkToken)
}

func NewVocabFromSlice[T constraints.Integer](lines []string, separator, unkToken string) (*Vocab[T], error) {
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
		vocab:    vocab,
		invVocab: invVocab,
		unkToken: unkToken,
	}, nil
}

func (v *Vocab[T]) Vocab() map[string]T {
	return v.vocab
}

func (v *Vocab[T]) UnkToken() string {
	return v.unkToken
}

func (v *Vocab[T]) TokenToID(token string) (T, error) {
	id, ok := v.vocab[token]
	if !ok {
		id, ok = v.vocab[v.unkToken]
		if !ok {
			return 0, fmt.Errorf("missing unknown token %q in vocab", v.unkToken)
		}
	}
	return id, nil
}

func (v *Vocab[T]) IDToToken(id T) (string, error) {
	token, ok := v.invVocab[id]
	if !ok {
		if _, ok = v.vocab[v.unkToken]; !ok {
			return "", fmt.Errorf("missing unknown token %q in vocab", v.unkToken)
		}
		token = v.unkToken
	}
	return token, nil
}

func (v *Vocab[T]) TokensToIDs(tokens []string) (ids []T, err error) {
	for _, token := range tokens {
		var id T
		id, err = v.TokenToID(token)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return
}

func (v *Vocab[T]) IDsToTokens(ids []T) (tokens []string, err error) {
	for _, id := range ids {
		var token string
		token, err = v.IDToToken(id)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	return
}

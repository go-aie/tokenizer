package tokenizer_test

import (
	"testing"

	aietokenizer "github.com/go-aie/tokenizer"
	"github.com/google/go-cmp/cmp"
	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/normalizer"
	"github.com/sugarme/tokenizer/pretokenizer"
)

func TestTokenizer_EncodeBatchTexts(t *testing.T) {
	lines := []string{
		"738	以",
		"1040	冬",
		"1282	及",
		"1914	天",
		"2519	开",
		"2763	想",
		"2834	我",
		"3266	春",
		"3907	法",
		"4853	的",
		"5166	秋",
		"5975	花",
		"6182	落",
		"7596	阳",
		"7794	风",
		"8000	OOV",
	}
	tk, err := newTokenizer(lines)
	if err != nil {
		t.Fatalf("err: %v\n", err)
	}

	tests := []struct {
		inTexts    []string
		wantIDs    [][]int
		wantTokens [][]string
	}{
		{
			inTexts: []string{"春天的花开秋天的风以及冬天的落阳", "我的想法"},
			wantIDs: [][]int{
				{3266, 1914, 4853, 5975, 2519, 5166, 1914, 4853, 7794, 738, 1282, 1040, 1914, 4853, 6182, 7596},
				{2834, 4853, 2763, 3907, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			wantTokens: [][]string{
				{"春", "天", "的", "花", "开", "秋", "天", "的", "风", "以", "及", "冬", "天", "的", "落", "阳"},
				{"我", "的", "想", "法", "", "", "", "", "", "", "", "", "", "", "", ""},
			},
		},
	}
	for _, tt := range tests {
		encodings, err := tk.EncodeBatchTexts(tt.inTexts, false)
		if err != nil {
			t.Errorf("err: %v\n", err)
		}

		var gotIDs [][]int
		for _, e := range encodings {
			gotIDs = append(gotIDs, e.Ids)
		}
		if !cmp.Equal(gotIDs, tt.wantIDs) {
			diff := cmp.Diff(gotIDs, tt.wantIDs)
			t.Errorf("Want - Got: %s", diff)
		}

		var gotTokens [][]string
		for _, e := range encodings {
			gotTokens = append(gotTokens, e.Tokens)
		}
		if !cmp.Equal(gotTokens, tt.wantTokens) {
			diff := cmp.Diff(gotTokens, tt.wantTokens)
			t.Errorf("Want - Got: %s", diff)
		}
	}
}

func newTokenizer(lines []string) (*aietokenizer.Tokenizer, error) {
	vocab, err := aietokenizer.NewVocabFromSlice[int](lines, "\t")
	if err != nil {
		return nil, err
	}
	m := aietokenizer.NewWordLevel(vocab.Vocab, "OOV")

	paddingStrategy := tokenizer.NewPaddingStrategy()
	paddingParams := tokenizer.PaddingParams{
		Strategy:  *paddingStrategy,
		Direction: tokenizer.Right, // padding right
	}
	tk := tokenizer.NewTokenizer(m)
	tk.WithPadding(&paddingParams)
	tk.WithNormalizer(normalizer.NewBertNormalizer(false, false, true, false)) // Handle Chinese chars
	tk.WithPreTokenizer(pretokenizer.NewBertPreTokenizer())

	return &aietokenizer.Tokenizer{Tokenizer: tk}, nil
}

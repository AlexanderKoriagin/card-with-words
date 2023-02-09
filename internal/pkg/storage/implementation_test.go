package storage

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"unicode"
)

func TestWords_Card(t *testing.T) {
	testData := []struct {
		lang Language
		qty  int
	}{
		{Russian, -1},
		{Russian, 0},
		{Russian, 1},
		{Russian, 8},
		{English, -2},
		{English, 0},
		{English, 1},
		{English, 8},
	}

	w := &words{}

	for _, td := range testData {
		card := w.Card(td.lang, td.qty)

		if td.qty <= 0 {
			assert.Equal(t, "", card)
			continue
		}

		separate := strings.Split(card, "\n")
		assert.Equal(t, td.qty, len(separate))

		for _, word := range card {
			if word != '\n' {
				assert.True(t, unicode.IsLetter(word))
			}
		}
	}

}

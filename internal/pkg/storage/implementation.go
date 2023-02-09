package storage

import (
	"cardWithWords/internal/pkg/data/words/english"
	"cardWithWords/internal/pkg/data/words/russian"
	"cardWithWords/internal/pkg/random"
	"strings"
)

type Language int

const (
	Russian Language = iota
	English
)

// Features interface to get "card" with qty words separated by \n
type Features interface {
	Card(lang Language, qty int) string
}

type words struct{}

func (w *words) Card(lang Language, qty int) string {
	var (
		str    string
		source []string
	)

	if qty <= 0 {
		return str
	}

	switch lang {
	case English:
		source = english.English
	default:
		source = russian.Russian
	}

	for qty > 0 {
		str += strings.ToUpper(source[random.Random(0, len(source))]) + "\n"
		qty--
	}

	return str[:len(str)-1]
}

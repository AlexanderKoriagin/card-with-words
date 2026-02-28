package groq

import (
	"strings"

	"cardWithWords/internal/pkg/base"

	"github.com/hashicorp/golang-lru/v2"
)

type cache struct {
	childRu *lru.Cache[int, string]
	childEn *lru.Cache[int, string]
	teenRu  *lru.Cache[int, string]
	teenEn  *lru.Cache[int, string]
	adultRu *lru.Cache[int, string]
	adultEn *lru.Cache[int, string]
}

func newCache(size int) (*cache, error) {
	childRu, err := lru.New[int, string](size)
	if err != nil {
		return nil, err
	}

	childEn, err := lru.New[int, string](size)
	if err != nil {
		return nil, err
	}

	teenRu, err := lru.New[int, string](size)
	if err != nil {
		return nil, err
	}

	teenEn, err := lru.New[int, string](size)
	if err != nil {
		return nil, err
	}

	adultRu, err := lru.New[int, string](size)
	if err != nil {
		return nil, err
	}

	adultEn, err := lru.New[int, string](size)
	if err != nil {
		return nil, err
	}

	return &cache{
		childRu: childRu,
		childEn: childEn,
		teenRu:  teenRu,
		teenEn:  teenEn,
		adultRu: adultRu,
		adultEn: adultEn,
	}, nil
}

func (c *cache) get(language base.Language, difficulty base.Difficulty) string {
	var values []string

	switch difficulty {
	case base.Child:
		switch language {
		case base.English:
			values = c.childEn.Values()
		default:
			values = c.childRu.Values()
		}
	case base.Teen:
		switch language {
		case base.English:
			values = c.teenEn.Values()
		default:
			values = c.teenRu.Values()
		}
	case base.Adult:
		switch language {
		case base.English:
			values = c.adultEn.Values()
		default:
			values = c.adultRu.Values()
		}
	}

	return strings.Join(values, base.SeparatorComma)
}

func (c *cache) add(language base.Language, difficulty base.Difficulty, words []string) {
	for i, word := range words {
		switch difficulty {
		case base.Child:
			switch language {
			case base.English:
				c.childEn.Add(i, word)
			default:
				c.childRu.Add(i, word)
			}
		case base.Teen:
			switch language {
			case base.English:
				c.teenEn.Add(i, word)
			default:
				c.teenRu.Add(i, word)
			}
		case base.Adult:
			switch language {
			case base.English:
				c.adultEn.Add(i, word)
			default:
				c.adultRu.Add(i, word)
			}
		}
	}
}

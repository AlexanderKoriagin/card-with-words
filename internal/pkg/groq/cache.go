package groq

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"cardWithWords/internal/pkg/base"

	"github.com/hashicorp/golang-lru/v2"
)

type cache struct {
	childRu *lru.Cache[string, string]
	childEn *lru.Cache[string, string]
	teenRu  *lru.Cache[string, string]
	teenEn  *lru.Cache[string, string]
	adultRu *lru.Cache[string, string]
	adultEn *lru.Cache[string, string]
}

func newCache(size int) (*cache, error) {
	childRu, err := lru.New[string, string](size)
	if err != nil {
		return nil, err
	}

	childEn, err := lru.New[string, string](size)
	if err != nil {
		return nil, err
	}

	teenRu, err := lru.New[string, string](size)
	if err != nil {
		return nil, err
	}

	teenEn, err := lru.New[string, string](size)
	if err != nil {
		return nil, err
	}

	adultRu, err := lru.New[string, string](size)
	if err != nil {
		return nil, err
	}

	adultEn, err := lru.New[string, string](size)
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
	for _, word := range words {
		fmt.Printf("Adding word %s to cache language %s difficulty %s\n", word, string(language), string(difficulty))
		hash := fmt.Sprintf("%s", sha256.Sum256([]byte(word)))
		switch difficulty {
		case base.Child:
			switch language {
			case base.English:
				c.childEn.Add(hash, word)
			default:
				c.childRu.Add(hash, word)
			}
		case base.Teen:
			switch language {
			case base.English:
				c.teenEn.Add(hash, word)
			default:
				c.teenRu.Add(hash, word)
			}
		case base.Adult:
			switch language {
			case base.English:
				c.adultEn.Add(hash, word)
			default:
				c.adultRu.Add(hash, word)
			}
		}
	}
}

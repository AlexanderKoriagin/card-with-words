package groq

import "cardWithWords/internal/pkg/base"

type Words interface {
	Card8Words(language base.Language, difficulty base.Difficulty) (*string, error)
}

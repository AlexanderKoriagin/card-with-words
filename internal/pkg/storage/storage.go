package storage

import "cardWithWords/internal/pkg/base"

type Words struct {
	Features Features
}

func Init() *Words {
	return &Words{Features: &words{}}
}

func (w *Words) GetRussian() string {
	return w.Features.Card(base.Russian, base.DefaultQty)
}

func (w *Words) GetEnglish() string {
	return w.Features.Card(base.English, base.DefaultQty)
}

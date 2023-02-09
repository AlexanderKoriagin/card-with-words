package storage

type Words struct {
	Features Features
}

func Init() *Words {
	return &Words{Features: &words{}}
}

func (w *Words) GetRussian() string {
	return w.Features.Card(Russian, 8)
}

func (w *Words) GetEnglish() string {
	return w.Features.Card(English, 8)
}

package telegram

import (
	"sync"

	"cardWithWords/internal/pkg/groq"
	"cardWithWords/internal/services"
)

type Telegram struct {
	Features Features
}

func Init(
	storage services.Data,
	groqWords groq.Words,
	wg *sync.WaitGroup,
	cDone chan struct{},
	cErr chan error,
) *Telegram {
	return &Telegram{
		Features: &Bot{
			WordsLocal: storage,
			WordsGroq:  groqWords,
			Wg:         wg,
			Channels: Channels{
				Done:   cDone,
				Errors: cErr,
			},
		},
	}
}

func (t *Telegram) PlayCards(token string) error {
	if err := t.Features.Open(token); err != nil {
		return err
	}

	go t.Features.Worker()
	return nil
}

package telegram

import (
	"cardWithWords/internal/services"
	"sync"
)

type Telegram struct {
	Features Features
}

func Init(storage services.Data, wg *sync.WaitGroup, cDone chan struct{}, cErr chan error) *Telegram {
	return &Telegram{
		Features: &Bot{
			Words: storage,
			Wg:    wg,
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

package telegram

import (
	"cardWithWords/internal/services"
	"fmt"
	tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sync"
)

const (
	buttonRus = "получить карту"
	buttonEng = "get a card"

	msgDefault = "нажми на кнопку / press the button"
)

type Channels struct {
	Done   chan struct{}
	Errors chan error
}

type Bot struct {
	Api      *tba.BotAPI
	Words    services.Data
	Wg       *sync.WaitGroup
	Channels Channels
}

type Features interface {
	Open(token string) error
	Worker()
}

func (b *Bot) Open(token string) (err error) {
	b.Api, err = tba.NewBotAPI(token)
	if err != nil {
		return err
	}

	return nil
}

// Worker generate words by request from user
func (b *Bot) Worker() {
	// Create a new UpdateConfig struct with an offset of 0. Offsets are used
	// to make sure Telegram knows we've handled previous values, and we don't
	// need them repeated.
	uc := tba.NewUpdate(0)

	// Tell Telegram we should wait up to 30 seconds on each request for an
	// update. This way we can get information just as quickly as making many
	// frequent requests without having to send nearly as many.
	uc.Timeout = 30

	// Start polling Telegram for updates
	updates := b.Api.GetUpdatesChan(uc)

	var keyboard = tba.NewReplyKeyboard(
		tba.NewKeyboardButtonRow(
			tba.NewKeyboardButton(buttonRus),
			tba.NewKeyboardButton(buttonEng),
		),
	)

	for {
		select {
		case <-b.Channels.Done:
			b.Wg.Done()
			return
		case u := <-updates:
			var card string

			if u.Message == nil {
				continue
			} else {
				switch u.Message.Text {
				case buttonRus:
					card = b.Words.GetRussian()
				case buttonEng:
					card = b.Words.GetEnglish()
				default:
					card = msgDefault
				}
			}

			msg := tba.NewMessage(u.Message.Chat.ID, card)
			msg.ReplyMarkup = keyboard
			//msg.ReplyMarkup = tba.NewReplyKeyboard(
			//	[]tba.KeyboardButton{
			//		{Text: buttonRus},
			//		{Text: buttonEng},
			//	},
			//)

			_, err := b.Api.Send(msg)
			if err != nil {
				b.Channels.Errors <- fmt.Errorf("[Worker] couldn't send msg to %s: %v\n", u.Message.From.UserName, err)
			}
		}
	}
}

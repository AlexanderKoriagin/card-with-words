package telegram

import (
	"fmt"
	"sync"

	"cardWithWords/internal/pkg/base"
	"cardWithWords/internal/pkg/groq"
	"cardWithWords/internal/services"

	tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	buttonRus         = "получить карту"
	buttonEng         = "get a card"
	buttonGroqChildRu = "GroqChildRu"
	buttonGroqTeenRu  = "GroqTeenRu"
	buttonGroqAdultRu = "GroqAdultRu"

	msgDefault = "нажми на кнопку / press the button"
)

type Channels struct {
	Done   chan struct{}
	Errors chan error
}

type Bot struct {
	Api        *tba.BotAPI
	WordsLocal services.Data
	WordsGroq  groq.Words
	Wg         *sync.WaitGroup
	Channels   Channels
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
			tba.NewKeyboardButton(buttonGroqChildRu),
			tba.NewKeyboardButton(buttonGroqTeenRu),
			tba.NewKeyboardButton(buttonGroqAdultRu),
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
					card = b.WordsLocal.GetRussian()
				case buttonEng:
					card = b.WordsLocal.GetEnglish()
				case buttonGroqChildRu:
					groqCard, err := b.WordsGroq.Card8Words(base.Russian, base.Child)
					if err != nil {
						b.Channels.Errors <- fmt.Errorf("[Worker] couldn't get card from Groq: %v\n", err)
						card = msgDefault
					}
					card = *groqCard
				case buttonGroqTeenRu:
					groqCard, err := b.WordsGroq.Card8Words(base.Russian, base.Teen)
					if err != nil {
						b.Channels.Errors <- fmt.Errorf("[Worker] couldn't get card from Groq: %v\n", err)
						card = msgDefault
					}
					card = *groqCard
				case buttonGroqAdultRu:
					groqCard, err := b.WordsGroq.Card8Words(base.Russian, base.Adult)
					if err != nil {
						b.Channels.Errors <- fmt.Errorf("[Worker] couldn't get card from Groq: %v\n", err)
						card = msgDefault
					}
					card = *groqCard
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
